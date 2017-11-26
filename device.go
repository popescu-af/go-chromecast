package cast

import "net"
import "fmt"

type Device struct {
	IP         net.IP
	Port       int
	Properties map[string]string
}

func (d Device) Addr() string {
	return fmt.Sprintf("%s:%d", d.IP, d.Port)
}

func (d Device) Name() string {
	return d.Properties["fn"]
}

func (d Device) ID() string {
	return d.Properties["id"]
}

func (d Device) Type() string {
	return d.Properties["md"]
}

func (d Device) Status() string {
	return d.Properties["rs"]
}
