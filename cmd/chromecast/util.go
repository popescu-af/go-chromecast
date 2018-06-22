package main

import (
	"context"
	"fmt"

	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/client"
	"github.com/oliverpool/go-chromecast/command"
	"github.com/oliverpool/go-chromecast/command/heartbeat"
	"github.com/oliverpool/go-chromecast/discovery"
	"github.com/oliverpool/go-chromecast/discovery/zeroconf"
	"github.com/oliverpool/go-chromecast/gogoprotobuf"
	"github.com/oliverpool/go-chromecast/net"
)

func GetClientWithStatus(ctx context.Context, logger chromecast.Logger) (chromecast.Client, chromecast.Status, error) {
	// Find device
	fmt.Print("Searching device... ")
	chr, err := GetDevice(ctx, logger)
	if err != nil {
		return nil, chromecast.Status{}, err
	}
	fmt.Println(chr.Addr() + " OK")

	// Connect client
	fmt.Print("Connecting client... ")
	client, err := ConnectedClient(ctx, chr.Addr(), logger)
	if err != nil {
		return nil, chromecast.Status{}, fmt.Errorf("could not connect to client: %v", err)
	}
	fmt.Println(" OK")

	launcher := command.Launcher{Requester: client}

	// Get receiver status
	fmt.Print("\nGetting receiver status...")
	status, err := launcher.Status()
	if err != nil {
		return nil, chromecast.Status{}, fmt.Errorf("could not get status: %v", err)
	}
	fmt.Println(" OK")
	fmt.Println(status.String())
	return client, status, nil
}

func GetDevice(ctx context.Context, logger chromecast.Logger) (*chromecast.Device, error) {
	// If IP is set, return device with corresponding IP
	if deviceIP != nil {
		return discovery.NewDevice(deviceIP, devicePort, nil), nil
	}

	// Otherwise search with matchers
	var matchers []discovery.DeviceMatcher
	if deviceName != "" {
		matchers = append(matchers, discovery.WithName(deviceName))
	}
	if deviceID != "" {
		matchers = append(matchers, discovery.WithID(deviceID))
	}
	chr, err := discovery.Service{Scanner: zeroconf.Scanner{Logger: logger}}.First(ctx, matchers...)
	if err != nil || chr == nil {
		return nil, fmt.Errorf("could not find a device: %v", err)
	}
	return chr, nil
}

// ConnectedClient will create a client and keep it connected
func ConnectedClient(ctx context.Context, addr string, logger chromecast.Logger) (*client.Client, error) {
	conn, err := net.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}

	serializer := gogoprotobuf.Serializer{
		Conn:   conn,
		Logger: logger,
	}
	c := client.New(&serializer, logger)

	go func() {
		<-ctx.Done()
		command.Close.Send(c)
		conn.Close()
	}()
	go heartbeat.RespondToPing(c)

	return c, command.Connect.Send(c)
}
