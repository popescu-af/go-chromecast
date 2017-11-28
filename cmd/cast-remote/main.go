package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oliverpool/go-chromecast"

	kitlog "github.com/go-kit/kit/log"

	"github.com/oliverpool/go-chromecast/cli"
	"github.com/oliverpool/go-chromecast/command"
	"github.com/oliverpool/go-chromecast/command/media"
)

var logger = kitlog.NewNopLogger()

func init() {
	logger = cli.NewLogger(os.Stdout)
	log.SetOutput(kitlog.NewStdlibAdapter(logger))
}

func fatalf(format string, a ...interface{}) int {
	fmt.Printf(format, a...)
	fmt.Println()
	return 1
}

func secondsToDuration(s float64) time.Duration {
	return time.Duration(s * float64(time.Second))
}

func main() {
	os.Exit(remote())
}

type localStatus struct {
	playing   bool
	volume    float64
	muted     bool
	time      time.Duration
	orderSent time.Time
}

func (cs *localStatus) update(cstatus chromecast.Status, mstatus media.Status) {
	if time.Since(cs.orderSent) < time.Second {
		return
	}
	fmt.Println("[up]")
	if cstatus.Volume != nil {
		cs.volume = *cstatus.Volume.Level
		cs.muted = *cstatus.Volume.Muted
	}
	cs.playing = (mstatus.PlayerState == "PLAYING")
	cs.time = time.Duration(mstatus.CurrentTime * float64(time.Second))
}

func (cs *localStatus) sent() {
	cs.orderSent = time.Now()
}

func remote() int {
	kill := make(chan struct{})
	ch := make(chan cli.KeyPress, 10)

	defer cli.ReadStdinKeys(ch, kill)()
	defer close(kill)

	fmt.Println("Ready:")

	var lstatus localStatus
	go func() {
		for {
			lstatus.update(chromecast.Status{}, media.Status{})
			time.Sleep(500 * time.Millisecond)
		}
	}()

	if true {
		for c := range ch {
			switch {
			case c.Type == cli.Escape:
				fmt.Println("bye")
				return 0
			case c.Type == cli.SpaceBar:
				fmt.Println("play/pause")
			case c.Type == cli.LowerCaseLetter:
				switch c.Key {
				case 'q':
					fmt.Println("stop")
					return 0
				case 'm':
					fmt.Println("mute/unmute")
				default:
					fmt.Println("key: " + string(c.Key))
				}
			case c.Type == cli.Arrow:
				switch c.Key {
				case cli.Up:
					fmt.Println("volume up")
				case cli.Down:
					fmt.Println("volume down")
				case cli.Left:
					fmt.Println("back 5s")
				case cli.Right:
					fmt.Println("forward 10s")
				}
			default:
				fmt.Println(c)
			}
			lstatus.sent()
		}
	}

	ctx := context.Background()

	fmt.Print("Searching device...")
	chr, err := cli.GetDevice(ctx, "", 0, "")
	if err != nil {
		return fatalf("could not get a device: %v", err)
	}
	fmt.Printf(" OK\n  '%s' (%s:%d)\n", chr.Name(), chr.IP, chr.Port)

	client, err := cli.NewClient(ctx, chr.Addr(), logger)
	if err != nil {
		return fatalf("could not get a client: %v", err)
	}

	fmt.Print("\nGetting receiver status...")
	status, err := command.Status.Get(client)
	if err != nil {
		return fatalf("could not get status: %v", err)
	}
	fmt.Println(" OK")
	cli.FprintStatus(os.Stdout, status)

	// Get Media app
	fmt.Print("\nLooking for a media app...")
	app, err := media.FromStatus(client, status)
	if err != nil {
		return fatalf(" not found")
	}
	fmt.Println(" OK")

	go app.UpdateStatus()
	_, err = app.GetStatus()
	if err != nil {
		return fatalf("could not get media status: %v", err)
	}
	session, err := app.CurrentSession()
	if err != nil {
		return fatalf("could not get a session")
	}

	if true {
		ch, err := session.Seek(media.Seek(15 * time.Second))
		if err != nil {
			return fatalf("could not pause")
		}
		fmt.Println(string(<-ch))
	}
	return 0
}
