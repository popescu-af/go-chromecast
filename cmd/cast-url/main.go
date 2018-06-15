package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/oliverpool/go-chromecast/command/media/defaultreceiver"
	"github.com/oliverpool/go-chromecast/command/media/defaultreceiver/tatort"
	"github.com/oliverpool/go-chromecast/command/media/defaultreceiver/tvnow"
	"github.com/oliverpool/go-chromecast/command/media/youtube"
	"github.com/oliverpool/go-chromecast/log"

	"github.com/oliverpool/go-chromecast/cli"
	"github.com/oliverpool/go-chromecast/command/media"
)

func fatalf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	fmt.Println()
	os.Exit(1)
}

func main() {
	ctx := context.Background()

	rawurl := "https://youtu.be/b-GIBLX3nAk"
	if len(os.Args) > 1 {
		rawurl = os.Args[1]
	}

	logger := log.NopLogger()
	if os.Getenv("DEBUG") == "" {
		logger = log.New(os.Stdout)
	}

	client, status, err := cli.FirstClientWithStatus(ctx, logger)
	if err != nil {
		fatalf(err.Error())
	}

	loaders := []struct {
		name   string
		loader media.URLLoader
	}{
		{"tatort", tatort.URLLoader},
		{"tvnow", tvnow.URLLoader},
		{"youtube", youtube.URLLoader},
		{"default", defaultreceiver.URLLoader},
	}

	for _, l := range loaders {
		loader, err := l.loader(rawurl)
		if err != nil {
			logger.Log("loader", l.name, "err", err)
			continue
		}
		c, err := loader(client, status)
		if err != nil {
			logger.Log("loader", l.name, "state", "loading", "err", err)
			continue
		}
		select {
		case <-c:
		case <-time.After(10 * time.Second):
			logger.Log("loader", l.name, "err", "load request didn't return after 10s")
		}
		return
	}
	fatalf("No supported loader found")
}
