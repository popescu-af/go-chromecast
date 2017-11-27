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

func (c command) SendTo(sender chromecast.Sender) error {
	return sender.Send(c.Envelope, c.Payload)
}
