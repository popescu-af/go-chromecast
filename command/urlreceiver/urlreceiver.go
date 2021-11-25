package urlreceiver

import (
	"fmt"
	"net/url"

	"github.com/popescu-af/go-chromecast"
	"github.com/popescu-af/go-chromecast/command"
	"github.com/popescu-af/go-chromecast/command/media"
)

// ID from https://github.com/DeMille/url-cast-receiver
const ID = "5CB45E5A"
const Namespace = "urn:x-cast:com.url.cast"

type App struct {
	*command.App
}

func LaunchAndConnect(client chromecast.Client, statuses ...chromecast.Status) (App, error) {
	// ignore statuses
	a, err := command.LaunchAndConnect(client, ID)
	if err != nil {
		return App{}, err
	}
	a.Envelope.Namespace = Namespace
	return App{
		App: a,
	}, nil
}

// Option to customize the loading
type Option func(command.Map)

func UseIframe(c command.Map) {
	c["type"] = "iframe"
}

func (a App) Load(url string, options ...media.Option) (<-chan []byte, error) {
	payload := command.Map{
		"type": "loc",
		"url":  url,
	}
	for _, opt := range options {
		opt(payload)
	}
	return a.Client.Request(a.Envelope, payload)
}

func URLLoader(rawurl string, options ...media.Option) (func(client chromecast.Client, statuses ...chromecast.Status) (<-chan []byte, error), error) {
	url, err := ExtractID(rawurl)
	if err != nil {
		return nil, err
	}
	return func(client chromecast.Client, statuses ...chromecast.Status) (<-chan []byte, error) {
		app, err := LaunchAndConnect(client, statuses...)
		if err != nil {
			return nil, err
		}
		return app.Load(url, options...)
	}, nil
}

func ExtractID(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf("could not parse url '%s': %v", rawurl, err)
	}
	return u.String(), nil
}
