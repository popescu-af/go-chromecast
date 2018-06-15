package discovery

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/oliverpool/go-chromecast"
)

// NewDevice returns an new chromecast device
func NewDevice(ip net.IP, port int, properties []string) *chromecast.Device {
	return &chromecast.Device{
		IP:         ip,
		Port:       port,
		Properties: parseProperties(properties),
	}
}

// Scanner scans for chromecast and pushes them onto the results channel (eventually multiple times)
// It must close the results channel in the returned function (which should return when the ctx is done)
type Scanner interface {
	Scan(ctx context.Context, results chan<- *chromecast.Device) (func() error, error)
}

// Service allows to discover chromecast via multiple means
type Service struct {
	Scanner Scanner
}

// First returns the first chromecast that is discovered by the scanner (and matching all matchers - if any)
func (s Service) First(ctx context.Context, matchers ...DeviceMatcher) (*chromecast.Device, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel child-ctx when the right client has been found

	result := make(chan *chromecast.Device, 1)

	worker, err := s.Scanner.Scan(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("could not initiliaze scanner: %v", err)
	}
	match := matchAll(matchers...)
	go worker()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case device := <-result:
			if match(device) {
				return device, nil
			}
		}
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
