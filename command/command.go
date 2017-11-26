package command

import cast "github.com/oliverpool/go-chromecast"

type command struct {
	Envelope cast.Envelope
	Payload  interface{}
}

type identifiableCommand struct {
	Envelope cast.Envelope
	Payload  cast.IdentifiablePayload
}

type requestFunc func(cast.Envelope, cast.IdentifiablePayload) (<-chan []byte, error)

type sendFunc func(cast.Envelope, interface{}) error

func (c command) SendTo(sender sendFunc) error {
	return sender(c.Envelope, c.Payload)
}
