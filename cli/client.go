package cli

import (
	"context"

	"github.com/oliverpool/go-chromecast"

	"github.com/oliverpool/go-chromecast/client"
	"github.com/oliverpool/go-chromecast/command"
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

	return c, command.Connect.Send(c)
}
