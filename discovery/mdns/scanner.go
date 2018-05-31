package mdns

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/mdns"
	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/discovery"

	"context"
)

// Scanner backed by the hashicorp/mdns package
type Scanner struct {
	// The chromecasts have 'Timeout' time to reply to each probe.
	Timeout time.Duration
}

// Scan repeatedly scans the network  and synchronously sends the chromecast found into the results channel.
// It finishes when the context is done.
func (s Scanner) Scan(ctx context.Context, results chan<- *chromecast.Device) (func() error, error) {
	return func() error {
		defer close(results)

		// generate entries
		entries := make(chan *mdns.ServiceEntry, 10)
		go func() {
			defer close(entries)
			for {
				if ctx.Err() != nil {
					return
				}
				mdns.Query(&mdns.QueryParam{
					Service: "_googlecast._tcp",
					Domain:  "local",
					Timeout: s.Timeout,
					Entries: entries,
				})
			}
		}()

		// decode entries
		for e := range entries {
			c, err := s.decode(e)
			if err != nil {
				continue
			}
			select {
			case results <- c:
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return ctx.Err()
	}, nil
}

// decode turns an mdns.ServiceEntry into a chromecast.Device
func (s Scanner) decode(entry *mdns.ServiceEntry) (*chromecast.Device, error) {
	if !strings.Contains(entry.Name, "._googlecast") {
		return nil, fmt.Errorf("fdqn '%s does not contain '._googlecast'", entry.Name)
	}

	var ip net.IP
	if len(entry.AddrV6) > 0 {
		ip = entry.AddrV6
	} else if len(entry.AddrV4) > 0 {
		ip = entry.AddrV4
	}

	return discovery.NewDevice(
		ip,
		entry.Port,
		entry.InfoFields,
	), nil
}
