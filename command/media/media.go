package media

import (
	"github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command"
)

// ID chosen from https://gist.github.com/jloutsenhizer/8855258
const ID = "CC1AD845"

type App struct {
	Envelope chromecast.Envelope
	Client   chromecast.Client
}

func New(client chromecast.Client) (a App, err error) {
	env, err := command.Launch.App(client, ID)
	if err != nil {
		return a, err
	}
	a.Envelope = env
	a.Client = client

	return a, command.Connect.SendTo(client, env.Destination)
}

func FromStatus(client chromecast.Client, st chromecast.Status) (a App, err error) {
	env, err := command.AppEnvFromStatus(st, ID, command.Status.Envelope.Source)
	if err != nil {
		return a, err
	}
	a.Envelope = env
	a.Client = client

	return a, command.Connect.SendTo(client, env.Destination)
}
