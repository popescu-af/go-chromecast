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
		logger, ctx, cancel := flags()
		defer cancel()

		client, status, err := GetClientWithStatus(ctx, logger)
		if err != nil {
			return fmt.Errorf("could not get a client: %w", err)
		}
		defer client.Close()

		return remote(ctx, cancel, logger, client, status)
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

func remote(
	initCtx context.Context,
	initCancel context.CancelFunc,
	logger chromecast.Logger,
	client chromecast.Client,
	status chromecast.Status,
) error {
	clientCtx := context.Background()
	clientCtx, clientCancel := context.WithCancel(clientCtx)

	cancel := func() {
		clientCancel()
		initCancel()
	}
	defer cancel()

	// Get Media app
	app, err := getMediaApp(client, status)
	if err != nil {
		return fmt.Errorf("could not get a media app: %w", err)
	}

	lstatus := local.New(status)
	// lstatus.UpdateMedia(app.LatestStatus()[0])

	var wg sync.WaitGroup
	wg.Add(1)

	var sessionFound uint32
	hasSession := func() bool {
		return atomic.LoadUint32(&sessionFound) == 1
	}

	session := new(media.Session)
	// var session *media.Session

	go func() {
		ch := make(chan cli.KeyPress, 10)
		go cli.ReadStdinKeyPresses(clientCtx, ch)

		defer cancel()
		defer wg.Done()

		processKeyInputs(
			ch,
			hasSession,
			session,
			lstatus,
			logger,
			command.Launcher{Requester: client}.AmpController(),
		)
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
			return fmt.Errorf("status could not be fetch: %w", err)
		}
	}
	fmt.Println(" OK")

	fmt.Print("Getting a session...")
	cs, err := app.CurrentSession()
	if err != nil {
		return fmt.Errorf("could not get a session: %w", err)
	}
	*session = *cs
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

func processKeyInputs(ch chan cli.KeyPress, hasSession func() bool, session *media.Session, lstatus *local.Status, logger chromecast.Logger, amp chromecast.AmpController) {

	forwardFactor := newStreakFactor()
	backwardFactor := newStreakFactor()

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
				amp.Quit()
				return
			case 'm':
				amp.Mute(lstatus.ToggleMute())
			default:
				logger.Log("msg", "unsupported lowercase", "key", string(c.Key), "type", c.Type)
			}
		case c.Type == cli.Arrow:
			switch c.Key {
			case cli.Up:
				amp.SetVolume(lstatus.IncrVolume(.1))
			case cli.Down:
				amp.SetVolume(lstatus.IncrVolume(-.1))
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
}

func getMediaApp(client chromecast.Client, status chromecast.Status) (app *media.App, err error) {
	fmt.Print("\nWaiting for a media app...")
	for {
		app, err = media.ConnectFromStatus(client, status)
		switch err {
		case nil:
			fmt.Println(" OK")
			go app.UpdateStatus()
			return app, nil
		case chromecast.ErrAppNotFound:
			time.Sleep(time.Second)
			fmt.Print(".")
			status, err = command.Launcher{Requester: client}.Status()
			if err != nil {
				return nil, fmt.Errorf("could not get status: %w", err)
			}
		default:
			return nil, fmt.Errorf("unexpected error: %w", err)
		}
	}
}
