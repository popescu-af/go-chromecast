package cli

import (
	"context"
	"fmt"

	"github.com/oliverpool/go-chromecast"

	"github.com/oliverpool/go-chromecast/client"
	"github.com/oliverpool/go-chromecast/command"
	"github.com/oliverpool/go-chromecast/command/heartbeat"
	"github.com/oliverpool/go-chromecast/gogoprotobuf"
	"github.com/oliverpool/go-chromecast/net"
)

// NewClient will send a Connect command
func NewClient(ctx context.Context, addr string, logger chromecast.Logger) (*client.Client, error) {
	conn, err := net.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}

	serializer := gogoprotobuf.Serializer{
		Conn:   conn,
		Logger: logger,
	}
	c := client.New(ctx, &serializer, logger)

	go func() {
		<-ctx.Done()
		command.Close.Send(c)
		conn.Close()
	}()
	go heartbeat.RespondToPing(c, c)

	return c, command.Connect.Send(c)
}

func GetClient(ctx context.Context, host string, port int, name string, logger chromecast.Logger) (*client.Client, error) {
	chr, err := GetDevice(ctx, host, port, name)
	if err != nil {
		return nil, fmt.Errorf("could not get a device: %v", err)
	}

	return NewClient(ctx, chr.Addr(), logger)
}
