package chromecast

import (
	"context"
)

type Scanner interface {
	// Scan scans for chromecast and pushes them onto the results channel (eventually multiple times)
	// It must close the results channel before returning when the ctx is done
	Scan(ctx context.Context, results chan<- *Device) error
}

type Envelope struct {
	Source, Destination, Namespace string
}

type Serializer interface {
	Receive() (Envelope, []byte, error)
	Send(Envelope, []byte) error
}

type IdentifiablePayload interface {
	SetRequestID(uint32)
}

type PayloadWithID struct {
	Type      string  `json:"type"`
	RequestID *uint32 `json:"requestId,omitempty"`
}

func (p *PayloadWithID) SetRequestID(id uint32) {
	p.RequestID = &id
}

type Client interface {
	Listen(env Envelope, responseType string, ch chan<- []byte)
	Send(env Envelope, payload interface{}) error
	Request(env Envelope, payload IdentifiablePayload) (<-chan []byte, error)
	Dispatch() error
	Close() error
}
