package cli

import (
	"context"

	"github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/discover"
	"github.com/oliverpool/go-chromecast/zeroconf"
)

func Scan(ctx context.Context, logger chromecast.Logger) chan *chromecast.Device {
	all := make(chan *chromecast.Device, 5)
	scanner := zeroconf.Scanner{
		Logger: logger,
	}
	go scanner.Scan(ctx, all)

	uniq := make(chan *chromecast.Device, 5)
	go discover.Uniq(all, uniq)
	return uniq
}
