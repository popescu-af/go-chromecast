package protocol

import (
	"crypto/tls"
	"fmt"
	"net"

	"golang.org/x/net/context"
)

func Dial(ctx context.Context, host net.IP, port int) (*tls.Conn, error) {
	deadline, _ := ctx.Deadline()
	dialer := &net.Dialer{
		Deadline: deadline,
	}
	return tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:%d", host, port), &tls.Config{
		InsecureSkipVerify: true,
	})
}
