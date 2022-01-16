package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/popescu-af/go-chromecast"

	"github.com/popescu-af/go-chromecast/command/media"
	"github.com/popescu-af/go-chromecast/command/media/defaultreceiver"
	"github.com/popescu-af/go-chromecast/command/media/defaultreceiver/tatort"
	"github.com/popescu-af/go-chromecast/command/media/defaultreceiver/tvnow"
	defaultvimeo "github.com/popescu-af/go-chromecast/command/media/defaultreceiver/vimeo"
	"github.com/popescu-af/go-chromecast/command/media/vimeo"
	"github.com/popescu-af/go-chromecast/command/media/youtube"
	"github.com/popescu-af/go-chromecast/command/urlreceiver"
	"github.com/spf13/cobra"
)

var loadRequestTimeout time.Duration

var useLoader string
var loaders = []namedLoader{
	{"tatort", tatort.URLLoader},
	{"tvnow", tvnow.URLLoader},
	{"vimeo", vimeo.URLLoader},
	{"youtube", youtube.URLLoader},
	{"default.vimeo", defaultvimeo.URLLoader},
	{"default", defaultreceiver.URLLoader},
	{"urlreceiver", urlreceiver.URLLoader},
}

var controlAfterwards bool

type namedLoader struct {
	name   string
	loader media.URLLoader
}

func (nl namedLoader) load(client chromecast.Client, status chromecast.Status, rawurl string) (<-chan []byte, error) {
	loader, err := nl.loader(rawurl)
	if err != nil {
		return nil, err
	}
	return loader(client, status)
}

func init() {
	loadCmd.Flags().DurationVarP(&loadRequestTimeout, "request-timeout", "r", 10*time.Second, "Duration to wait for a reply to the load request")
	var ll []string
	for _, l := range loaders {
		ll = append(ll, l.name)
	}
	loadCmd.Flags().StringVarP(&useLoader, "loader", "l", "", "Loader to use (supported loaders: "+strings.Join(ll, ", ")+")")
	loadCmd.Flags().BoolVarP(&controlAfterwards, "control", "c", false, "Launch control afterwards")
	rootCmd.AddCommand(loadCmd)
}

var loadCmd = &cobra.Command{
	Use:   "load [url]",
	Short: "Load a URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rawurl := strings.TrimSpace(args[0])

		logger, ctx, cancel := flags(cmd)
		defer cancel()
		return loadURL(ctx, cancel, rawurl, useLoader, controlAfterwards, logger)
	},
}

func loadURL(ctx context.Context, cancel context.CancelFunc, rawurl, useLoader string, controlAfterwards bool, logger chromecast.Logger) error {
	client, status, err := GetClientWithStatus(ctx, logger)
	if err != nil {
		return fmt.Errorf("could not get a client: %w", err)
	}
	defer client.Close()

	for _, l := range loaders {
		var c <-chan []byte
		var err error

		if useLoader != "" {
			if l.name != useLoader {
				continue
			}
			c, err = l.load(client, status, rawurl)
			if err != nil {
				return err
			}
		} else {
			c, err = l.load(client, status, rawurl)
			if err != nil {
				logger.Log("loader", l.name, "state", "loading", "err", err)
				continue
			}
			fmt.Printf("Loading with %s\n", l.name)
		}
		select {
		case <-c:
		case <-time.After(loadRequestTimeout):
			logger.Log("loader", l.name, "err", "load request didn't return after 10s")
		}
		if controlAfterwards {
			return remote(ctx, cancel, logger, client, status)
		}
		return nil
	}
	if useLoader != "" {
		var ll []string
		for _, l := range loaders {
			ll = append(ll, l.name)
		}
		return fmt.Errorf("unknown loader '%s' (supported loaders: %s)", useLoader, strings.Join(ll, ", "))
	}
	return fmt.Errorf("no supported loader found for %s", rawurl)
}
