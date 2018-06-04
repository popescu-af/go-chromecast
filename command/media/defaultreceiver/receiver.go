package defaultreceiver

import (
	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command/media"
)

const ID = "CC1AD845"

func LaunchAndConnect(client chromecast.Client, statuses ...chromecast.Status) (*media.App, error) {
	return media.LaunchAndConnect(client, ID, statuses...)
}
