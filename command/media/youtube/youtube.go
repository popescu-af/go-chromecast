package youtube

import (
	"fmt"
	"net/url"
	"path"

	"github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command/media"
)

const ID = "233637DE"

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
		ContentType: "x-youtube/video",
		StreamType:  "BUFFERED",
	}
	return a.App.Load(item, options...)
}

func (a App) LoadURL(rawurl string, options ...media.Option) (<-chan []byte, error) {
	id, err := ExtractID(rawurl)
	if err != nil {
		return nil, err
	}
	return a.Load(id, options...)
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
		"youtube-nocookie.com":     struct{}{},
		"www.youtube-nocookie.com": struct{}{},
		"youtu.be":                 struct{}{},
		"youtube.com":              struct{}{},
		"www.youtube.com":          struct{}{},
	}
	if _, ok := hosts[u.Host]; !ok {
		return "", fmt.Errorf("unsupported host: %s", u.Host)
	}
	if id := u.Query().Get("v"); id != "" {
		return id, nil
	}
	if id := path.Base(u.Path); id != "" {
		return id, nil
	}
	return "", fmt.Errorf("could not find id inside URL: %s", rawurl)
}
