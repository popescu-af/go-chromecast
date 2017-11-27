package command

import "github.com/oliverpool/go-chromecast"

type Command struct {
	Envelope chromecast.Envelope
	Payload  interface{}
}

func (c Command) Send(sender chromecast.Sender) error {
	return sender.Send(c.Envelope, c.Payload)
}

func (c Command) SendTo(sender chromecast.Sender, destination string) error {
	c.Envelope.Destination = destination
	return c.Send(sender)
}

type identifiableCommand struct {
	Envelope chromecast.Envelope
	Payload  chromecast.IdentifiablePayload
}

type Map map[string]interface{}

func (m Map) SetRequestID(ID uint32) {
	m["requestId"] = ID
}
