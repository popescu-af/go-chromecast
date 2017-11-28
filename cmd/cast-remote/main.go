package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oliverpool/go-chromecast/cli/local"
	"github.com/oliverpool/go-chromecast/command/volume"

	"github.com/gosuri/uiprogress"

	kitlog "github.com/go-kit/kit/log"

	"github.com/oliverpool/go-chromecast/cli"
	"github.com/oliverpool/go-chromecast/command"
	"github.com/oliverpool/go-chromecast/command/media"
)

var logger = kitlog.NewNopLogger()

func init() {
	// logger = cli.NewLogger(os.Stdout)
	log.SetOutput(kitlog.NewStdlibAdapter(logger))
}

func fatalf(format string, a ...interface{}) int {
	fmt.Printf(format, a...)
	fmt.Println()
	return 1
}

func main() {
	os.Exit(remote())
}

func newStreakFactor() func() int64 {
	var streakStart time.Time
	var previousHit time.Time
	return func() int64 {
		now := time.Now()
		defer func() { previousHit = now }()
		if now.Sub(previousHit) < 50*time.Millisecond {
			switch {
			case now.Sub(streakStart) > 3*time.Second:
				return 12
			case now.Sub(streakStart) > 2*time.Second:
				return 8
			case now.Sub(streakStart) > time.Second:
				return 2
			}
		} else {
			streakStart = now
		}
		return 1
	}
}

func remote() int {
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
		return fatalf(" not found: %v", err)
	}
	fmt.Println(" OK")

	go app.UpdateStatus()

	fmt.Print("Looking for a playing item...")
	appStatus, err := app.GetStatus()
	for err != nil || len(appStatus) == 0 || appStatus[0].Item == nil || appStatus[0].Item.Duration == 0 {
		if ctx.Err() != nil {
			return fatalf("could not get media status: %v", err)
		}
		appStatus, err = app.GetStatus()
	}
	fmt.Println(" OK")

	fmt.Print("Getting a session...")
	session, err := app.CurrentSession()
	if err != nil {
		return fatalf("could not get a session: %v", err)
	}
	fmt.Println(" OK\n")

	kill := make(chan struct{})
	ch := make(chan cli.KeyPress, 10)

	defer cli.ReadStdinKeys(ch, kill)()
	defer close(kill)

	total := int(appStatus[0].Item.Duration)

	bar := uiprogress.AddBar(total)
	bar.Width = 40
	uiprogress.Start()

	lstatus := local.New(status)
	lstatus.UpdateMedia(app.LatestStatus()[0])

	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return lstatus.PlayerState()
	})
	bar.AppendFunc(func(b *uiprogress.Bar) string {
		return lstatus.TimeStatus()
	})

	forwardFactor := newStreakFactor()
	backwardFactor := newStreakFactor()

	go func() {
		for {
			app.GetStatus()
			elapsed := lstatus.UpdateMedia(app.LatestStatus()[0])
			bar.Set(elapsed)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	for c := range ch {
		switch {
		case c.Type == cli.Escape:
			uiprogress.Stop()
			fmt.Println("bye")
			return 0
		case c.Type == cli.SpaceBar:
			if lstatus.TogglePlay() {
				session.Play()
			} else {
				session.Pause()
			}
		case c.Type == cli.LowerCaseLetter:
			switch c.Key {
			case 'q':
				uiprogress.Stop()
				fmt.Println("stop")
				ch, _ := session.Stop()
				<-ch
				return 0
			case 'm':
				volume.Mute(client, lstatus.ToggleMute())
				// default:
				// 	fmt.Println("key: " + string(c.Key))
			}
		case c.Type == cli.Arrow:
			switch c.Key {
			case cli.Up:
				volume.Set(client, lstatus.IncrVolume(.1))
			case cli.Down:
				volume.Set(client, lstatus.IncrVolume(-.1))
			case cli.Left:
				diff := -time.Duration(backwardFactor()) * 5 * time.Second
				session.Seek(media.Seek(lstatus.SeekBy(diff)))
			case cli.Right:
				diff := time.Duration(forwardFactor()) * 10 * time.Second
				session.Seek(media.Seek(lstatus.SeekBy(diff)))
			}
			// default:
			// 	fmt.Println(c)
		}
	}

	return 0
}
