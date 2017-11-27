package command

import (
	chromecast "github.com/oliverpool/go-chromecast"
)

type volumeController struct {
	Envelope chromecast.Envelope
}

func (v volumeController) Get(requester chromecast.Requester) (*chromecast.Volume, error) {
	st, err := Status.Get(requester)
	return st.Volume, err
}

func (v volumeController) Set(requester chromecast.Requester, level float64) (<-chan []byte, error) {
	vol := chromecast.Volume{
		Level: &level,
	}
	payload := Map{
		"type":   "SET_VOLUME",
		"volume": vol,
	}
	return requester.Request(v.Envelope, payload)
}

func (v volumeController) Mute(requester chromecast.Requester, muted bool) (<-chan []byte, error) {
	vol := chromecast.Volume{
		Muted: &muted,
	}
	payload := Map{
		"type":   "SET_VOLUME",
		"volume": vol,
	}
	return requester.Request(v.Envelope, payload)
}

var Volume = volumeController{
	Envelope: Status.Envelope,
}

/*
var Volume = volumeController{
	Envelope: connectionEnv,
	Payload:  Map{"type": "CLOSE"},
}
*/
