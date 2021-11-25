package client

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/popescu-af/go-chromecast"
)

func New(serializer chromecast.Serializer, logger chromecast.Logger) *Client {
	c := Client{
		Serializer: serializer,
		Logger:     logger,
	}

	go func() {
		var lastErr string
		nbErr := 0
		for nbErr <= 5 {
			if err := c.Dispatch(); err != nil {
				logger.Log("step", "dispatch", "err", err)
				if err.Error() == lastErr {
					nbErr++
				} else {
					lastErr = err.Error()
					nbErr = 1
				}
			} else {
				nbErr = 0
			}
		}
		logger.Log("step", "dispatch-abort", "err", fmt.Errorf("same error %d times: %s", nbErr, lastErr))
	}()

	return &c
}

type Client struct {
	chromecast.Serializer
	Logger     chromecast.Logger
	AfterClose []func()

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
		c.forwardResponse(*payID.RequestID, pay)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if env.Destination == "*" {
		// broadcast
		for _, listeners := range c.listeners {
			nonBlockingForwardTo(listeners, payID.Type, pay)
		}
	} else if listeners, ok := c.listeners[env]; ok {
		nonBlockingForwardTo(listeners, payID.Type, pay)
	}
	return err
}

func nonBlockingForwardTo(listeners map[string][]chan<- []byte, key string, payload []byte) {
	if typeListeners, ok := listeners[key]; ok {
		for _, ch := range typeListeners {
			// non blocking send
			select {
			case ch <- payload:
			default:
			}
		}
	}
}

func (c *Client) forwardResponse(id uint32, pay []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if requester, ok := c.pending[id]; ok {
		requester <- pay
		close(requester)
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
	for _, cb := range c.AfterClose {
		cb()
	}
	return nil
}
