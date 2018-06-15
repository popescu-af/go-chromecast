package main

import (
	"context"
	"fmt"
	"os"

	"github.com/oliverpool/go-chromecast/cli"
	"github.com/oliverpool/go-chromecast/command/media"
	"github.com/oliverpool/go-chromecast/log"
)

func fatalf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	fmt.Println()
	os.Exit(1)
}

func main() {
	ctx := context.Background()

	logger := log.NopLogger()
	if os.Getenv("DEBUG") == "" {
		logger = log.New(os.Stdout)
	}

	client, status, err := cli.FirstClientWithStatus(ctx, logger)
	if err != nil {
		fatalf(err.Error())
	}

	// Get media app
	fmt.Print("\nLooking for a media app...")
	app, err := media.ConnectFromStatus(client, status)
	if err != nil {
		fatalf(" nothing found: %v", err)
	}
	fmt.Println(" OK")

	// Get media app status
	fmt.Print("Getting media app status...")
	st, err := app.Status()
	if err != nil {
		fatalf("could not get media status: %v", err)
	}
	fmt.Println(" OK")
	for _, s := range st {
		if s.Item != nil {
			fmt.Printf("  Item: %s\n", s.Item.ContentId)
			fmt.Printf("  Type: %s\n", s.Item.ContentType)
			fmt.Printf("  Stream: %s\n", s.Item.StreamType)
			fmt.Printf("  Duration: %s\n", s.Item.Duration)
			fmt.Printf("  Metadata: %#v\n", s.Item.Metadata)
		}
		fmt.Printf("  Current Time: %s\n", s.CurrentTime)
		fmt.Printf("  State: %s\n", s.PlayerState)
		fmt.Printf("  Rate: %.2f\n", s.PlaybackRate)
	}
}
