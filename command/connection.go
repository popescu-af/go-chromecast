package command

import (
	"github.com/oliverpool/go-chromecast"
)

const DefaultSource = "sender-0"
const DefaultDestination = "receiver-0"

var connectionEnv = chromecast.Envelope{
	Source:      DefaultSource,
	Destination: DefaultDestination,
	Namespace:   "urn:x-cast:com.google.cast.tp.connection",
}

var Connect = Command{
	Envelope: connectionEnv,
	Payload:  Type("CONNECT"),
}

var Close = Command{
	Envelope: connectionEnv,
	Payload:  Type("CLOSE"),
}
