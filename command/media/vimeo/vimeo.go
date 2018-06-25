package vimeo

import (
	"fmt"
	"net/url"
	"path"

	"github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command/media"
)

const ID = "7742C69E"

type App struct {
	*media.App
}

func LaunchAndConnect(client chromecast.Client, statuses ...chromecast.Status) (App, error) {
	app, err := media.LaunchAndConnect(client, ID, statuses...)
	return App{app}, err
}

func (a App) Load(id string, options ...media.Option) (<-chan []byte, error) {
	item := media.Item{
		ContentID:   id,
		ContentType: "application/dash+xml",
		StreamType:  "BUFFERED",
	}
	return a.App.Load(item, options...)
}

func URLLoader(rawurl string, options ...media.Option) (func(client chromecast.Client, statuses ...chromecast.Status) (<-chan []byte, error), error) {
	id, err := ExtractID(rawurl)
	if err != nil {
		return nil, err
	}
	return func(client chromecast.Client, statuses ...chromecast.Status) (<-chan []byte, error) {
		app, err := LaunchAndConnect(client, statuses...)
		if err != nil {
			return nil, err
		}
		return app.Load(id, options...)
	}, nil
}

func ExtractID(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf("could not parse url '%s': %v", rawurl, err)
	}

	hosts := map[string]struct{}{
		"vimeo.com": struct{}{},
	}
	if _, ok := hosts[u.Host]; !ok {
		return "", fmt.Errorf("unsupported host: %s", u.Host)
	}
	if id := path.Base(u.Path); id != "" {
		return "/videos/" + id, nil
	}
	return "", fmt.Errorf("could not find id inside URL: %s", rawurl)
}
