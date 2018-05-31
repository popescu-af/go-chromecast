package discovery_test

import (
	"fmt"
	"strings"
	"testing"

	"context"

	"github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/discovery"
	"github.com/oliverpool/go-chromecast/mock"
)

func TestFirstDirect(t *testing.T) {
	scan := mock.Scanner{
		ScanFunc: func(ctx context.Context, results chan<- *chromecast.Device) (func() error, error) {
			return func() error {
				results <- &chromecast.Device{}
				close(results)
				return nil
			}, nil
		},
	}

	service := discovery.Service{Scanner: &scan}

	ctx := context.Background()

	first, err := service.First(ctx)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if first == nil {
		t.Errorf("a client should have been found")
	}
	if scan.ScanFuncCalled != 1 {
		t.Errorf("scanner should have been called once, and not %d times", scan.ScanFuncCalled)
	}
}

func TestFirstCancelled(t *testing.T) {
	scan := mock.Scanner{
		ScanFunc: func(ctx context.Context, results chan<- *chromecast.Device) (func() error, error) {
			return func() error {
				<-ctx.Done()
				return nil
			}, nil
		},
	}

	service := discovery.Service{Scanner: &scan}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	first, err := service.First(ctx)
	if err != ctx.Err() {
		t.Errorf("unexpected error %v", err)
	}
	if first != nil {
		t.Errorf("a client should not have been found")
	}
	if scan.ScanFuncCalled > 1 {
		t.Errorf("scanner should have been called at most once, and not %d times", scan.ScanFuncCalled)
	}
}

func TestNamedDirect(t *testing.T) {
	scan := mock.Scanner{}
	done := make(chan struct{})
	scan.ScanFunc = func(ctx context.Context, results chan<- *chromecast.Device) (func() error, error) {
		return func() error {
			defer close(results)
			results <- &chromecast.Device{}
			c := &chromecast.Device{
				Properties: map[string]string{
					"fn": "casti",
				},
			}
			results <- c
			results <- &chromecast.Device{}
			select {
			case results <- &chromecast.Device{}:
				t.Error("channel should have been full")
			case <-ctx.Done():
			}
			close(done)
			return nil
		}, nil
	}

	service := discovery.Service{Scanner: &scan}

	ctx := context.Background()

	first, err := service.Named(ctx, "casti")
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if first == nil {
		t.Errorf("a client should have been found")
	}
	if first.Name() != "casti" {
		t.Errorf("the client should been named 'casti' and not '%s'", first.Name())
	}
	if scan.ScanFuncCalled != 1 {
		t.Errorf("scanner should have been called once, and not %d times", scan.ScanFuncCalled)
	}
	<-done
}

func TestNamedCancelled(t *testing.T) {
	scan := mock.Scanner{}
	done := make(chan struct{})
	scan.ScanFunc = func(ctx context.Context, results chan<- *chromecast.Device) (func() error, error) {
		return func() error {
			defer close(results)
			for {
				select {
				case results <- &chromecast.Device{}:
				case <-ctx.Done():
					close(done)
					return nil
				}
			}
		}, nil
	}

	service := discovery.Service{Scanner: &scan}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	first, err := service.Named(ctx, "casti")
	if err != ctx.Err() {
		t.Errorf("unexpected error %v", err)
	}
	if err != ctx.Err() {
		t.Errorf("unexpected error %v", err)
	}
	if first != nil {
		t.Errorf("a client should not have been found")
	}
	if scan.ScanFuncCalled > 1 {
		t.Errorf("scanner should have been called at most once, and not %d times", scan.ScanFuncCalled)
	}
	<-done
}

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

func TestNewDevice(t *testing.T) {

	txt := `id=87cf98a003f1f1dbd2efe6d19055a617|ve=04|md=Chromecast|ic=/setup/icon.png|fn=Chromecast PO|ca=5|st=0|bs=FA8FCA7EE8A9|rs=`

	exp := map[string]string{
		"id": "87cf98a003f1f1dbd2efe6d19055a617",
		"ve": "04",
		"md": "Chromecast",
		"ic": "/setup/icon.png",
		"fn": "Chromecast PO",
		"ca": "5",
		"st": "0",
		"bs": "FA8FCA7EE8A9",
		"rs": "",
	}

	chr := discovery.NewDevice(nil, 0, strings.Split(txt, "|"))

	res := chr.Properties
	if !mapEqual(exp, res) {
		t.Errorf("expected %s; found %s", exp, res)
	}
}

func mapEqual(m1, m2 map[string]string) bool {
	if m1 == nil {
		return m2 == nil
	}
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			fmt.Println(k, v1, v2, ok)
			return false
		}
	}
	return true
}
