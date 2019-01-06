package defaultreceiver

import (
	"fmt"
	"net/url"
	"path"

	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command/media"
)

const ID = "CC1AD845"

func LaunchAndConnect(client chromecast.Client, statuses ...chromecast.Status) (*media.App, error) {
	return media.LaunchAndConnect(client, ID, statuses...)
}

func URLLoader(rawurl string, options ...media.Option) (func(client chromecast.Client, statuses ...chromecast.Status) (<-chan []byte, error), error) {
	contentType, err := ExtractType(rawurl)
	if err != nil {
		return nil, err
	}
	return func(client chromecast.Client, statuses ...chromecast.Status) (<-chan []byte, error) {
		app, err := LaunchAndConnect(client, statuses...)
		if err != nil {
			return nil, err
		}
		return app.Load(media.Item{
			ContentID:   rawurl,
			ContentType: contentType,
			StreamType:  "BUFFERED",
		}, options...)
	}, nil
}

func ExtractType(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf("could not parse url '%s': %v", rawurl, err)
	}
	if t := u.Query().Get("ext"); t != "" {
		return contentTypeFromExtension(t), nil
	}
	t := contentTypeFromExtension(path.Ext(u.Path))
	if t == "" {
		return "", fmt.Errorf("could not find suitable content-type for '%s' (use the 'ext=.mpd' to force it)", path.Ext(u.Path))
	}
	return contentTypeFromExtension(path.Ext(u.Path)), nil
}

func contentTypeFromExtension(ext string) string {
	switch ext {
	case ".m3u8":
		return "application/x-mpegurl"
	case ".mpd":
		return "application/dash+xml"
	case ".ism":
		return "application/vnd.ms-sstr+xml"
	case ".mp4":
		return "video/mp4"
	default:
		return ""
	}
}
