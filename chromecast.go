package cast

import "net"

type Chromecast struct {
	IP   net.IP
	Port int
	Info map[string]string
}

func (c Chromecast) Name() string {
	return c.Info["fn"]
}

func (c Chromecast) ID() string {
	return c.Info["id"]
}

func (c Chromecast) Device() string {
	return c.Info["md"]
}

func (c Chromecast) Status() string {
	return c.Info["rs"]
}
