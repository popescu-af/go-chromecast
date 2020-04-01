package main

import (
	"context"
	"fmt"
	"net"

	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/discovery"
	"github.com/oliverpool/go-chromecast/discovery/zeroconf"
)

func init() {
	rootCmd.PersistentFlags().IPVar(&deviceFinder.IP, "ip", nil, "Specify chromecast IP")
	rootCmd.PersistentFlags().IntVar(&deviceFinder.Port, "port", 8009, "Specify chromecast port (ignored if IP is not set)")
	rootCmd.PersistentFlags().StringVarP(&deviceFinder.Name, "name", "n", "", "Specify chromecast name (ignored if IP is set)")
	rootCmd.PersistentFlags().StringVar(&deviceFinder.ID, "id", "", "Specify chromecast ID (ignored if IP is set)")
}

type deviceFinderConstraints struct {
	Name string
	ID   string
	IP   net.IP
	Port int
}

var deviceFinder deviceFinderConstraints

func (df deviceFinderConstraints) GetDevice(ctx context.Context, logger chromecast.Logger) (*chromecast.Device, error) {
	// If IP is set, return device with corresponding IP
	if df.IP != nil {
		return discovery.NewDevice(df.IP, df.Port, nil), nil
	}

	// Otherwise search with matchers
	var matchers []discovery.DeviceMatcher
	if df.Name != "" {
		matchers = append(matchers, discovery.WithName(df.Name))
	}
	if df.ID != "" {
		matchers = append(matchers, discovery.WithID(df.ID))
	}
	chr, err := discovery.Service{Scanner: zeroconf.Scanner{Logger: logger}}.First(ctx, matchers...)
	if err != nil || chr == nil {
		return nil, fmt.Errorf("could not find a device: %w", err)
	}
	return chr, nil
}
