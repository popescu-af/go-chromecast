package command

import (
	"github.com/oliverpool/go-chromecast"
)

var connectionEnv = chromecast.Envelope{
	Source:      "sender-0",
	Destination: "receiver-0",
	Namespace:   "urn:x-cast:com.google.cast.tp.connection",
}

var Connect = Command{
	Envelope: connectionEnv,
	Payload:  Map{"type": "CONNECT"},
}

var Close = Command{
	Envelope: connectionEnv,
	Payload:  Map{"type": "CLOSE"},
}
