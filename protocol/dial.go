package protocol

import (
	"crypto/tls"
	"net"

	"context"
)

func Dial(ctx context.Context, addr string) (*tls.Conn, error) {
	deadline, _ := ctx.Deadline()
	dialer := &net.Dialer{
		Deadline: deadline,
	}
	return tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
	})
}
