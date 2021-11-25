package cli

import (
	"context"
	"fmt"

	"github.com/popescu-af/go-chromecast"

	"github.com/popescu-af/go-chromecast/client"
	"github.com/popescu-af/go-chromecast/command"
	"github.com/popescu-af/go-chromecast/command/heartbeat"
	"github.com/popescu-af/go-chromecast/discovery"
	"github.com/popescu-af/go-chromecast/discovery/zeroconf"
	"github.com/popescu-af/go-chromecast/gogoprotobuf"
	"github.com/popescu-af/go-chromecast/net"
)

// FirstClientWithStatus find a device, connects a client and get its status (and is verbose)
func FirstClientWithStatus(ctx context.Context, logger chromecast.Logger, matchers ...discovery.DeviceMatcher) (chromecast.Client, chromecast.Status, error) {
	// Find device
	fmt.Print("Searching device... ")
	chr, err := discovery.Service{Scanner: zeroconf.Scanner{Logger: logger}}.First(ctx, matchers...)
	if err != nil || chr == nil {
		return nil, chromecast.Status{}, fmt.Errorf("could not discover a device: %w", err)
	}
	fmt.Println(chr.Addr() + " OK")

	// Connect client
	fmt.Print("Connecting client... ")
	client, err := ConnectedClient(ctx, chr.Addr(), logger)
	if err != nil {
		return nil, chromecast.Status{}, fmt.Errorf("could not connect to client: %w", err)
	}
	fmt.Println(" OK")

	launcher := command.Launcher{Requester: client}

	// Get receiver status
	fmt.Print("\nGetting receiver status...")
	status, err := launcher.Status()
	if err != nil {
		return nil, chromecast.Status{}, fmt.Errorf("could not get status: %w", err)
	}
	fmt.Println(" OK")
	fmt.Println(status.String())
	return client, status, nil
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
