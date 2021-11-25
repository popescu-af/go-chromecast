package defaultreceiver

import (
	"fmt"
	"net/url"
	"path"

	chromecast "github.com/popescu-af/go-chromecast"
	"github.com/popescu-af/go-chromecast/command/media"
)

const ID = "CC1AD845"

func LaunchAndConnect(client chromecast.Client, statuses ...chromecast.Status) (*media.App, error) {
	return media.LaunchAndConnect(client, ID, statuses...)
}

func URLLoader(rawurl string, tracks []media.Track, options ...media.Option) (func(client chromecast.Client, statuses ...chromecast.Status) (<-chan []byte, error), error) {
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
			//// EDIT by me !!!!
			Tracks: tracks,
			SubtitleTrackStyle: media.TextTrackStyle{
				BackgroundColor:   "#00000000",
				EdgeColor:         "#FFFFFFFF",
				EdgeType:          "OUTLINE",
				FontFamily:        "ARIAL",
				FontGenericFamily: "SANS_SERIF",
				FontScale:         1,
				FontStyle:         "NORMAL",
				ForegroundColor:   "#FF0000FF",
			},

			////
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
	if fragment := u.Fragment; fragment != "" {
		return contentTypeFromExtension(fragment), nil
	}
	t := contentTypeFromExtension(path.Ext(u.Path))
	if t == "" {
		return "", fmt.Errorf("could not find suitable content-type for '%s' (use the 'ext=.mpd' or '#.mp4' to force it)", path.Ext(u.Path))
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
