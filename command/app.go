package command

import (
	"fmt"

	chromecast "github.com/oliverpool/go-chromecast"
)

type App struct {
	Envelope chromecast.Envelope
	Client   chromecast.Client
}

func LaunchAndConnect(client chromecast.Client, id string, statuses ...chromecast.Status) (*App, error) {
	st, err := Launcher{Requester: client}.Launch(id, statuses...)
	if err != nil {
		return nil, fmt.Errorf("could not launch app: %w", err)
	}
	app := st.AppWithID(id)
	if app == nil {
		return nil, fmt.Errorf("the launched app could not be found")
	}
	if app.TransportId == nil {
		return nil, fmt.Errorf("the app has no transportId")
	}
	return ConnectTo(client, *app.TransportId)
}

func ConnectFromStatus(client chromecast.Client, st chromecast.Status, namespace string) (*App, error) {
	destination, err := st.FirstDestinationSupporting(namespace)
	if err != nil {
		return nil, err
	}
	a, err := ConnectTo(client, destination)
	if err != nil {
		return nil, err
	}
	a.Envelope.Namespace = namespace
	return a, nil
}

func ConnectTo(client chromecast.Client, destination string) (*App, error) {
	a := &App{
		Envelope: chromecast.Envelope{
			Source:      DefaultSource,
			Destination: destination,
		},
		Client: client,
	}
	return a, Connect.SendTo(client, destination)
}
