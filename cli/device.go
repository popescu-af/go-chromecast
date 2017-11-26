package cli

import (
	"net"
	"time"

	cast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/discover"
	"github.com/oliverpool/go-chromecast/mdns"
	"context"
)

// GetDevice will try to get a casting device.
// If host is set, it will be used (with its port).
// Otherwise, if name is set, a chromecast will be looked-up by name.
// Otherwise the first chromecast found will be returned.
func GetDevice(ctx context.Context, host string, port int, name string) (*cast.Device, error) {
	if host != "" {
		ips, err := net.LookupIP(host)
		if err != nil {
			return nil, err
		}
		return &cast.Device{
			IP:         ips[0],
			Port:       port,
			Properties: make(map[string]string),
		}, nil
	}

	find := discover.Service{
		Scanner: mdns.Scanner{
			Timeout: 3 * time.Second,
		},
	}
	if name != "" {
		return find.Named(ctx, name)
	}
	return find.First(ctx)
}
