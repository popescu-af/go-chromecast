package cli

import (
	"github.com/barnybug/go-cast/client"
	"github.com/barnybug/go-cast/command"
	"github.com/barnybug/go-cast/protocol"
	"golang.org/x/net/context"
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
