package volume

import (
	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command"
)

func Get(requester chromecast.Requester) (*chromecast.Volume, error) {
	st, err := command.Status.Get(requester)
	return st.Volume, err
}

func Set(requester chromecast.Requester, level float64) (<-chan []byte, error) {
	vol := chromecast.Volume{
		Level: &level,
	}
	env := command.Status.Envelope
	payload := command.Map{
		"type":   "SET_VOLUME",
		"volume": vol,
	}
	return requester.Request(env, payload)
}

func Mute(requester chromecast.Requester, muted bool) (<-chan []byte, error) {
	vol := chromecast.Volume{
		Muted: &muted,
	}
	env := command.Status.Envelope
	payload := command.Map{
		"type":   "SET_VOLUME",
		"volume": vol,
	}
	return requester.Request(env, payload)
}
