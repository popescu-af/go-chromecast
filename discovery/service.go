package discovery

import (
	"context"
	"net"
	"strings"

	"github.com/oliverpool/go-chromecast"
)

// Scanner scans for chromecast and pushes them onto the results channel (eventually multiple times)
// It must close the results channel before returning when the ctx is done
type Scanner interface {
	Scan(ctx context.Context, results chan<- *chromecast.Device) error
}

// Service allows to discover chromecast via multiple means
type Service struct {
	Scanner chromecast.Scanner
}

// First returns the first chromecast that is discovered by the scanner
func (s Service) First(ctx context.Context) (*chromecast.Device, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel child-ctx when the first client has been found

	result := make(chan *chromecast.Device, 1)

	go s.Scanner.Scan(ctx, result)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case client := <-result:
		return client, nil
	}
}

// Named returns the first chromecast that is discovered by the scanner with the given name
func (s Service) Named(ctx context.Context, name string) (*chromecast.Device, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel child-ctx when the right client has been found

	result := make(chan *chromecast.Device, 1)

	go s.Scanner.Scan(ctx, result)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case client := <-result:
			if client.Name() == name {
				return client, nil
			}
		}
	}
}

// Uniq forward all client deduplicated
func Uniq(in <-chan *chromecast.Device, out chan<- *chromecast.Device) {
	seen := make(map[string]struct{})
	for c := range in {
		if c == nil {
			continue
		}
		if _, ok := seen[c.ID()]; ok {
			continue
		}
		out <- c
		seen[c.ID()] = struct{}{}
	}
	close(out)
}

// NewDevice returns an new chromecast device
func NewDevice(ip net.IP, port int, properties []string) *chromecast.Device {
	return &chromecast.Device{
		IP:         ip,
		Port:       port,
		Properties: parseProperties(properties),
	}
}

// parseProperties into a string map
// Input: {"key1=value1", "key2=value2"}
func parseProperties(s []string) map[string]string {
	m := make(map[string]string, len(s))
	for _, v := range s {
		s := strings.SplitN(v, "=", 2)
		if len(s) == 2 {
			m[s[0]] = s[1]
		}
	}
	return m
}
