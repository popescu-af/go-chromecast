package heartbeat

import (
	"github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command"
)

func RespondToPing(client chromecast.Client) {
	pingEnvelope := chromecast.Envelope{
		Source:      "Tr@n$p0rt-0",
		Destination: "Tr@n$p0rt-0",
		Namespace:   "urn:x-cast:com.google.cast.tp.heartbeat",
	}

	ch := make(chan []byte, 1)
	client.Listen(pingEnvelope, "PING", ch)

	for range ch {
		client.Send(pingEnvelope, command.Type("PONG"))
	}
}
