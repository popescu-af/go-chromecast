package command

import (
	"encoding/json"

	"github.com/popescu-af/go-chromecast"
)

type Command struct {
	Envelope chromecast.Envelope
	Payload  interface{}
}

type Type string

func (t Type) MarshalJSON() (b []byte, err error) {
	return json.Marshal(struct {
		Type string `json:"type"`
	}{Type: string(t)})
}

func (c Command) Send(sender chromecast.Sender) error {
	return sender.Send(c.Envelope, c.Payload)
}

func (c Command) SendTo(sender chromecast.Sender, destination string) error {
	c.Envelope.Destination = destination
	return c.Send(sender)
}

type Map map[string]interface{}

// SetRequestID helps fullfills the chromecast.IdentifiablePayload interface
func (m Map) SetRequestID(ID uint32) {
	m["requestId"] = ID
}
