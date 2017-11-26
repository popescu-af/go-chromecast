package command

import (
	"github.com/barnybug/go-cast"
)

var connectionEnv = cast.Envelope{
	Source:      "sender-0",
	Destination: "receiver-0",
	Namespace:   "urn:x-cast:com.google.cast.tp.connection",
}

var Connect = command{
	Envelope: connectionEnv,
	Payload:  cast.PayloadWithID{Type: "CONNECT"},
}

var Close = command{
	Envelope: connectionEnv,
	Payload:  cast.PayloadWithID{Type: "CLOSE"},
}
