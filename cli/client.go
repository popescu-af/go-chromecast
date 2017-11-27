package cli

import (
	"context"

	"github.com/oliverpool/go-chromecast/client"
	"github.com/oliverpool/go-chromecast/command"
	"github.com/oliverpool/go-chromecast/gogoprotobuf"
	"github.com/oliverpool/go-chromecast/net"
)

// NewClient will send a Connect command
func NewClient(ctx context.Context, addr string) (*client.Client, error) {
	conn, err := net.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}

	serializer := gogoprotobuf.Serializer{
		Conn: conn,
	}
	c := client.New(ctx, &serializer)

	go func() {
		<-ctx.Done()
		command.Close.SendTo(c)
		conn.Close()
	}()

	return c, command.Connect.SendTo(c)
}
