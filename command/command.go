package command

import "github.com/oliverpool/go-chromecast"

type command struct {
	Envelope chromecast.Envelope
	Payload  interface{}
}

type identifiableCommand struct {
	Envelope chromecast.Envelope
	Payload  chromecast.IdentifiablePayload
}

type requestFunc func(chromecast.Envelope, chromecast.IdentifiablePayload) (<-chan []byte, error)

type sendFunc func(chromecast.Envelope, interface{}) error

func (c command) SendTo(sender sendFunc) error {
	return sender(c.Envelope, c.Payload)
}
