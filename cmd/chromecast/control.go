package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/oliverpool/go-chromecast/streak"

	"github.com/gosuri/uiprogress"
	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/cli"
	"github.com/oliverpool/go-chromecast/cli/local"
	"github.com/oliverpool/go-chromecast/command"
	"github.com/oliverpool/go-chromecast/command/media"
	"github.com/spf13/cobra"
)

type Controller interface {
	Close() error

	// session
	Play() error
	Pause() error
	Seek(t time.Duration) error
	Stop() error

	// launcher
	Mute(muted bool) error
	SetVolume(level float64) error
	Quit() error
}

type Progress interface {
	Set(float64) error
	Close() error
}

func init() {
	rootCmd.AddCommand(controlCmd)
}

var controlCmd = &cobra.Command{
	Use:   "control",
	Short: "Control a chromecast",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return remote()
	},
}

func newStreakFactor() func() int64 {
	s := streak.New(50*time.Millisecond, streak.Factor{
		After: 3 * time.Second,
		Value: 6,
	}, streak.Factor{
		After: 2 * time.Second,
		Value: 4,
	}, streak.Factor{
		After: 1 * time.Second,
		Value: 2,
	})
	return s.UpdatedFactor
}

func remote() error {
	clientCtx := context.Background()
	clientCtx, clientCancel := context.WithCancel(clientCtx)

	logger, initCtx, initCancel := flags()
	cancel := func() {
		clientCancel()
		initCancel()
	}
	defer cancel()

	client, status, err := GetClientWithStatus(initCtx, logger)
	if err != nil {
		return fmt.Errorf("could not get a client: %v", err)
	}
	defer client.Close()
	launcher := command.Launcher{Requester: client}

	// Get Media app
	fmt.Print("\nWaiting for a media app...")
	var app *media.App
	for {
		app, err = media.ConnectFromStatus(client, status)
		if err == nil {
			fmt.Println(" OK")
			break
		}
		if clientCtx.Err() != nil {
			return fmt.Errorf("%v", clientCtx.Err())
		}
		if err == chromecast.ErrAppNotFound {
			select {
			case <-clientCtx.Done():
				return fmt.Errorf("interrupted: %v", clientCtx.Err())
			case <-time.After(time.Second):
			}
			fmt.Print(".")
			status, err = launcher.Status()
			if err != nil {
				return fmt.Errorf("could not get status: %v", err)
			}
			continue
		} else if err != nil {
			return fmt.Errorf(" failed: %v", err)
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
	for len(appStatus) == 0 || appStatus[0].Item == nil || appStatus[0].Item.Duration.Seconds() == 0 {
		select {
		case <-clientCtx.Done():
			return fmt.Errorf("interrupted: %v", clientCtx.Err())
		case <-time.After(time.Second):
		}
		fmt.Print(".")
		appStatus, err = app.Status()
		if err != nil {
			return fmt.Errorf("status could not be fetch: %v", err)
		}
	}
	fmt.Println(" OK")

	fmt.Print("Getting a session...")
	session, err = app.CurrentSession()
	if err != nil {
		return fmt.Errorf("could not get a session: %v", err)
	}
	fmt.Println(" OK")

	fmt.Println("\n Play/Pause: <space>  Seek: ←/→  Volume: ↑/↓/m  Stop: s  Quit: q  Disconnect: <Esc>")

	total := int(appStatus[0].Item.Duration.Seconds())

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
			if len(app.LatestStatus()) > 0 {
				elapsed := lstatus.UpdateMedia(app.LatestStatus()[0])
				bar.Set(elapsed)
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	wg.Wait()

	return nil
}
