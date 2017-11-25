package cast

import (
	"golang.org/x/net/context"
)

type Scanner interface {
	// Scan scans for chromecast and pushes them onto the results channel (eventually multiple times)
	// It must close the results channel before returning when the ctx is done
	Scan(ctx context.Context, results chan<- *Client) error
}

type Message struct {
	Header
	Payload
}
type Header struct {
	Type      string
	RequestID *int
}
type Payload []byte

type Serializer interface {
	Send(payload interface{}, sourceId, destinationId, namespace string) error
	Receive() (*Message, error)
}
