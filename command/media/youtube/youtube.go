package youtube

import (
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
