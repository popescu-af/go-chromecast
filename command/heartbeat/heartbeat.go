package heartbeat

import (
	"github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command"
)

var pingEnv = chromecast.Envelope{
	Source:      "Tr@n$p0rt-0",
	Destination: "Tr@n$p0rt-0",
	Namespace:   "urn:x-cast:com.google.cast.tp.heartbeat",
}

func RespondToPing(client chromecast.Client) {
	ch := make(chan []byte, 1)
	client.Listen(pingEnv, "PING", ch)

	for range ch {
		client.Send(pingEnv, command.Map{"type": "PONG"})
	}
}

func SendPing(listener chromecast.Listener, sender chromecast.Sender) {
	sender.Send(pingEnv, command.Map{"type": "PING"})
}
