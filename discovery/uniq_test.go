package discovery_test

import (
	"testing"

	chromecast "github.com/popescu-af/go-chromecast"
	"github.com/popescu-af/go-chromecast/discovery"
)

func TestUniq(t *testing.T) {
	in := make(chan *chromecast.Device, 10)
	in <- &chromecast.Device{}
	in <- &chromecast.Device{}
	in <- &chromecast.Device{}
	in <- &chromecast.Device{}
	c := &chromecast.Device{
		Properties: map[string]string{
			"id": "123",
		},
	}
	in <- c
	in <- c
	close(in)

	out := make(chan *chromecast.Device, 2)
	discovery.Uniq(in, out)
	c = <-out
	if c.ID() != "" {
		t.Errorf("unexpected ID: %s", c.ID())
	}
	c = <-out
	if c.ID() != "123" {
		t.Errorf("unexpected ID: %s", c.ID())
	}
	c, ok := <-out
	if ok {
		t.Error("out should have been closed")
	}
}
