package cli

import (
	"time"

	cast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/discover"
	"github.com/oliverpool/go-chromecast/mdns"
	"context"
)

func Scan(ctx context.Context) chan *cast.Device {
	all := make(chan *cast.Device, 5)
	scanner := mdns.Scanner{
		Timeout: 3 * time.Second,
	}
	go scanner.Scan(ctx, all)

	uniq := make(chan *cast.Device, 5)
	go discover.Uniq(all, uniq)
	return uniq
}
