package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/oliverpool/go-chromecast"

	"github.com/oliverpool/go-chromecast/cli/local"
	"github.com/oliverpool/go-chromecast/command"

	"github.com/gosuri/uiprogress"

	"github.com/oliverpool/go-chromecast/cli"
	"github.com/oliverpool/go-chromecast/command/media"
	"github.com/oliverpool/go-chromecast/command/receiver"
)

func fatalf(format string, a ...interface{}) int {
	fmt.Printf(format, a...)
	fmt.Println()
	return 1
}

func main() {
	ctx := context.Background()
	logger := cli.NewLogger(os.Stdout)

	var cancel context.CancelFunc
	if timeout, err := time.ParseDuration(os.Getenv("TIMEOUT")); err == nil {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		logger.Log("timeout", timeout)
		defer cancel()
	}
	os.Exit(remote(ctx, logger))
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
				return 6
			case now.Sub(streakStart) > 2*time.Second:
				return 4
			case now.Sub(streakStart) > time.Second:
				return 2
			}
		} else {
			streakStart = now
		}
		return 1
	}
}

func remote(ctx context.Context, logger chromecast.Logger) int {
	clientCtx := context.Background()
	clientCtx, clientCancel := context.WithCancel(clientCtx)

	ctx, initCancel := context.WithCancel(ctx)
	cancel := func() {
		clientCancel()
		initCancel()
	}
	defer cancel()

	client, status, err := cli.FirstClientWithStatus(ctx, logger)
	if err != nil {
		fatalf(err.Error())
	}
	launcher := receiver.Launcher{Requester: client}

	// Get Media app
	fmt.Print("\nWaiting for a media app...")
	var app *media.App
	for {
		app, err = media.FromStatus(client, status)
		if err == nil {
			fmt.Println(" OK")
			break
		}
		if ctx.Err() != nil {
			return fatalf("%v", ctx.Err())
		}
		if err == command.ErrAppNotFound {
			select {
			case <-ctx.Done():
				return fatalf("interrupted: %v", ctx.Err())
			case <-time.After(time.Second):
			}
			fmt.Print(".")
			status, err = launcher.Status()
			if err != nil {
				return fatalf("could not get status: %v", err)
			}
			continue
		} else if err != nil {
			return fatalf(" failed: %v", err)
		}
	}

	go app.UpdateStatus()

	kill := make(chan struct{})
	ch := make(chan cli.KeyPress, 10)

	defer cli.ReadStdinKeys(ch, kill)()
	defer close(kill)

	lstatus := local.New(status)
	// lstatus.UpdateMedia(app.LatestStatus()[0])

	forwardFactor := newStreakFactor()
	backwardFactor := newStreakFactor()

	var wg sync.WaitGroup
	wg.Add(1)

	var sessionFound uint32
	hasSession := func() bool {
		return atomic.LoadUint32(&sessionFound) == 1
	}

	var session *media.Session

	go func() {
		defer cancel()
		defer wg.Done()

		for c := range ch {
			switch {
			case c.Type == cli.Escape:
				if hasSession() {
					uiprogress.Stop()
					fmt.Println("bye")
				}
				return
			case c.Type == cli.SpaceBar && hasSession():
				if lstatus.TogglePlay() {
					session.Play()
				} else {
					session.Pause()
				}
			case c.Type == cli.LowerCaseLetter:
				switch c.Key {
				case 's':
					if !hasSession() {
						continue
					}
					uiprogress.Stop()
					fmt.Println("stop")
					ch, _ := session.Stop()
					<-ch
					return
				case 'q':
					if hasSession() {
						uiprogress.Stop()
					}
					fmt.Println("quit")
					launcher.Stop()
					return
				case 'm':
					launcher.Mute(lstatus.ToggleMute())
				default:
					logger.Log("msg", "unsupported lowercase", "key", string(c.Key), "type", c.Type)
				}
			case c.Type == cli.Arrow:
				switch c.Key {
				case cli.Up:
					launcher.SetVolume(lstatus.IncrVolume(.1))
				case cli.Down:
					launcher.SetVolume(lstatus.IncrVolume(-.1))
				case cli.Left:
					if !hasSession() {
						continue
					}
					diff := -time.Duration(backwardFactor()) * 5 * time.Second
					session.Seek(media.Seek(lstatus.SeekBy(diff)))
				case cli.Right:
					if !hasSession() {
						continue
					}
					diff := time.Duration(forwardFactor()) * 10 * time.Second
					session.Seek(media.Seek(lstatus.SeekBy(diff)))
				default:
					logger.Log("msg", "unsupported arrow", "key", c.Key, "type", c.Type)
				}
			default:
				logger.Log("msg", "unsupported key", "key", c.Key, "type", c.Type)
			}
		}
	}()

	// Get loaded item
	fmt.Print("Waiting for a loaded item...")
	appStatus := app.LatestStatus()
	for len(appStatus) == 0 || appStatus[0].Item == nil || appStatus[0].Item.Duration == 0 {
		select {
		case <-ctx.Done():
			return fatalf("interrupted: %v", ctx.Err())
		case <-time.After(time.Second):
		}
		fmt.Print(".")
		appStatus, err = app.Status()
		if err != nil {
			return fatalf("status could not be fetch: %v", err)
		}
	}
	fmt.Println(" OK")

	fmt.Print("Getting a session...")
	session, err = app.CurrentSession()
	if err != nil {
		return fatalf("could not get a session: %v", err)
	}
	fmt.Println(" OK")

	fmt.Println("\n Play/Pause: <space>  Seek: ←/→  Volume: ↑/↓/m  Stop: s  Quit: q  Disconnect: <Esc>")

	total := int(appStatus[0].Item.Duration)

	bar := uiprogress.AddBar(total)
	bar.Width = 40
	uiprogress.Start()

	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return lstatus.PlayerState()
	})
	bar.AppendFunc(func(b *uiprogress.Bar) string {
		return lstatus.TimeStatus()
	})

	atomic.StoreUint32(&sessionFound, 1)

	go func() {
		for {
			app.Status()
			elapsed := lstatus.UpdateMedia(app.LatestStatus()[0])
			bar.Set(elapsed)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	wg.Wait()

	return 0
}
