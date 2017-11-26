package cli

import (
	"github.com/oliverpool/go-chromecast/client"
	"github.com/oliverpool/go-chromecast/command"
	"github.com/oliverpool/go-chromecast/protocol"
	"context"
)

// NewClient will send a Connect command
func NewClient(ctx context.Context, addr string) (*client.Client, error) {
	conn, err := protocol.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}

	serializer := protocol.Serializer{
		Conn: conn,
	}
	c := client.New(ctx, &serializer)

	go func() {
		<-ctx.Done()
		command.Close.SendTo(c.Send)
		conn.Close()
	}()

	return c, command.Connect.SendTo(c.Send)
}
