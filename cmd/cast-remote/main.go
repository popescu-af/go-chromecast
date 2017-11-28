package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gosuri/uiprogress"
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
	// fmt.Println("[up]")
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
	if false {
		ctx := context.Background()

		fmt.Print("Searching device...")
		client, err := cli.GetClient(ctx, "", 0, "", logger)
		if err != nil {
			return fatalf("could not get a client: %v", err)
		}
		fmt.Println(" OK")

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
		_ = session
	}

	/*



	 */

	kill := make(chan struct{})
	ch := make(chan cli.KeyPress, 10)

	defer cli.ReadStdinKeys(ch, kill)()
	defer close(kill)

	fmt.Println("Ready:")

	duration := 159 * time.Minute
	total := int(duration.Seconds())

	bar := uiprogress.AddBar(total)
	bar.Width = 40
	uiprogress.Start()

	var lstatus localStatus
	lstatus.playing = true
	var fakeTime float64

	bar.PrependFunc(func(b *uiprogress.Bar) string {
		if lstatus.playing {
			return " Playing "
		}
		return "[paused] "
	})
	bar.AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%8s/%8s", lstatus.time.Round(time.Second), duration.Round(time.Second))
	})

	fakePlay := true
	go func() {
		for {
			state := "PLAYING"
			if !lstatus.playing {
				state = "PAUSED"
			}
			lstatus.update(chromecast.Status{}, media.Status{
				CurrentTime: fakeTime,
				PlayerState: state,
			})
			bar.Set(int(lstatus.time.Seconds()))
			time.Sleep(500 * time.Millisecond)
			if fakePlay {
				fakeTime += .5
			}
		}
	}()

	if true {
		for c := range ch {
			switch {
			case c.Type == cli.Escape:
				uiprogress.Stop()
				fmt.Println("bye")
				return 0
			case c.Type == cli.SpaceBar:
				lstatus.playing = !lstatus.playing
				fakePlay = lstatus.playing
			case c.Type == cli.LowerCaseLetter:
				switch c.Key {
				case 'q':
					uiprogress.Stop()
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
					lstatus.time -= 5 * time.Second
					fakeTime -= 5
				case cli.Right:
					lstatus.time += 10 * time.Second
					fakeTime += 10
				}
			default:
				fmt.Println(c)
			}
			lstatus.sent()
		}
	}

	return 0
}
