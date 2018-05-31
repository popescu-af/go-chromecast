package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	kitlog "github.com/go-kit/kit/log"

	"github.com/oliverpool/go-chromecast/cli"
	"github.com/oliverpool/go-chromecast/command/media"
	"github.com/oliverpool/go-chromecast/command/receiver"
	"github.com/oliverpool/go-chromecast/discovery"
	"github.com/oliverpool/go-chromecast/discovery/zeroconf"
)

var logger = kitlog.NewNopLogger()

func init() {
	if os.Getenv("DEBUG") != "" {
		logger = cli.NewLogger(os.Stdout)
	}
	log.SetOutput(kitlog.NewStdlibAdapter(logger))
}

func fatalf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	fmt.Println()
	os.Exit(1)
}

func secondsToDuration(s float64) time.Duration {
	return time.Duration(s * float64(time.Second))
}

func main() {
	ctx := context.Background()

	// Find device
	fmt.Print("Searching device... ")
	chr, err := discovery.Service{Scanner: zeroconf.Scanner{Logger: logger}}.First(ctx)
	if err != nil {
		fatalf("could not discover a device: %v", err)
	}
	fmt.Println(chr.Addr() + " OK")

	// Connect client
	fmt.Print("Connecting client... ")
	client, err := cli.ConnectedClient(ctx, chr.Addr(), logger)
	if err != nil {
		fatalf("could not connect to client: %v", err)
	}
	fmt.Println(" OK")

	launcher := receiver.Launcher{
		Requester: client,
	}

	// Get receiver status
	fmt.Print("\nGetting receiver status...")
	status, err := launcher.Status()
	if err != nil {
		fatalf("could not get status: %v", err)
	}
	fmt.Println(" OK")
	fmt.Println(status.String())

	// Get media app
	fmt.Print("\nLooking for a media app...")
	app, err := media.FromStatus(client, status)
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
			fmt.Printf("  Duration: %s\n", secondsToDuration(s.Item.Duration))
			fmt.Printf("  Metadata: %#v\n", s.Item.Metadata)
		}
		fmt.Printf("  Current Time: %s\n", secondsToDuration(s.CurrentTime))
		fmt.Printf("  State: %s\n", s.PlayerState)
		fmt.Printf("  Rate: %.2f\n", s.PlaybackRate)
	}
}
