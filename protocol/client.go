package protocol

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/barnybug/go-cast"
)

type Client struct {
	cast.Serializer
	SourceID      string
	DestinationID string
	Namespace     string

	requestID uint32
	mu        sync.Mutex
	pending   map[uint32]chan<- cast.Payload
	listeners map[string][]chan<- cast.Payload
}

func (c Client) shouldReceive(h cast.Header) bool {
	return h.DestinationID == "*" || (h.Namespace == c.Namespace &&
		h.SourceID == c.DestinationID && h.DestinationID == c.SourceID)
}

func (c *Client) Listen(responseType string, ch chan<- cast.Payload) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.listeners == nil {
		c.listeners = make(map[string][]chan<- cast.Payload, 1)
	}
	c.listeners[responseType] = append(c.listeners[responseType], ch)
}

func (c *Client) Send(payload interface{}) error {
	return c.Serializer.Send(payload, c.SourceID, c.DestinationID, c.Namespace)
}

func (c *Client) Request(payload cast.IdentifiablePayload) (<-chan cast.Payload, error) {
	id := atomic.AddUint32(&c.requestID, 1)

	payload.SetRequestID(id)
	response := make(chan cast.Payload, 1)
	c.mu.Lock()
	if c.pending == nil {
		c.pending = make(map[uint32]chan<- cast.Payload, 1)
	}
	c.pending[id] = response
	c.mu.Unlock()

	err := c.Send(payload)
	if err != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, err
	}
	return response, nil
}

func (c *Client) Dispatch() error {
	message, err := c.Serializer.Receive()
	if err != nil {
		return err
	}

	h := message.Header

	if !c.shouldReceive(h) {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Println("received", *message.RequestID)

	if message.RequestID != nil {
		if requester, ok := c.pending[*h.RequestID]; ok {
			requester <- message.Payload
			delete(c.pending, *h.RequestID)
		}
	}

	if listeners, ok := c.listeners[h.Type]; ok {
		for _, ch := range listeners {
			select {
			case ch <- message.Payload:
			default:
			}
		}
	}

	return nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, requester := range c.pending {
		close(requester)
		delete(c.pending, id)
	}
	for responseType, listeners := range c.listeners {
		for _, ch := range listeners {
			close(ch)
		}
		delete(c.listeners, responseType)
	}
	return c.Serializer.Close()
}
