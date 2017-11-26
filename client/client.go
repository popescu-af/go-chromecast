package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/log"
)

func New(ctx context.Context, serializer chromecast.Serializer) *Client {
	c := Client{
		Serializer: serializer,
	}

	go func() {
		log.Println("dispatching...")
		for ctx.Err() == nil {
			err := c.Dispatch()
			if err != nil {
				log.Println("dispatch failed:", err)
			}
		}
	}()

	return &c
}

type Client struct {
	chromecast.Serializer

	requestID uint32
	mu        sync.Mutex
	pending   map[uint32]chan<- []byte
	listeners map[chromecast.Envelope]map[string][]chan<- []byte
}

func (c *Client) Listen(env chromecast.Envelope, responseType string, ch chan<- []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.listeners == nil {
		c.listeners = make(map[chromecast.Envelope]map[string][]chan<- []byte, 1)
	}

	var ok bool
	var types map[string][]chan<- []byte

	if types, ok = c.listeners[env]; !ok {
		types = make(map[string][]chan<- []byte)
		c.listeners[env] = types
	}

	types[responseType] = append(types[responseType], ch)
}

func (c *Client) Send(env chromecast.Envelope, payload interface{}) error {
	pay, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %s", err)
	}
	return c.Serializer.Send(env, pay)
}

func (c *Client) Request(env chromecast.Envelope, payload chromecast.IdentifiablePayload) (<-chan []byte, error) {
	id := atomic.AddUint32(&c.requestID, 1)

	payload.SetRequestID(id)
	response := make(chan []byte, 1)
	c.mu.Lock()
	if c.pending == nil {
		c.pending = make(map[uint32]chan<- []byte, 1)
	}
	c.pending[id] = response
	c.mu.Unlock()

	err := c.Send(env, payload)
	if err != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, err
	}
	return response, nil
}

func (c *Client) Dispatch() error {
	env, pay, err := c.Serializer.Receive()
	if err != nil {
		return err
	}

	var payID chromecast.PayloadWithID

	err = json.Unmarshal(pay, &payID)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal into PayloadWithID: %s", err)
	}

	if payID.RequestID != nil {
		c.sendResponse(*payID.RequestID, pay)
	}

	if env.Namespace == "*" {
		// broadcast
		for _, envs := range c.listeners {
			for _, listeners := range envs {
				for _, ch := range listeners {
					// non blocking send
					select {
					case ch <- pay:
					default:
					}
				}
			}
		}
	} else if listeners, ok := c.listeners[env]; ok {
		if typeListeners, ok := listeners[payID.Type]; ok {
			for _, ch := range typeListeners {
				// non blocking send
				select {
				case ch <- pay:
				default:
				}
			}
		}
	}
	return nil
}

func (c *Client) sendResponse(id uint32, pay []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if requester, ok := c.pending[id]; ok {
		requester <- pay
		delete(c.pending, id)
	}
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, requester := range c.pending {
		close(requester)
		delete(c.pending, id)
	}
	for _, envs := range c.listeners {
		for responseType, listeners := range envs {
			for _, ch := range listeners {
				close(ch)
			}
			delete(envs, responseType)
		}
	}
	return nil
}
