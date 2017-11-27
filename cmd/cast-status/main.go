package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	fmt.Print("Searching device...")
	chr, err := cli.GetDevice(ctx, "", 0, "")
	if err != nil {
		fatalf("could not get a device: %v", err)
	}
	fmt.Printf(" OK\n  '%s' (%s:%d)\n", chr.Name(), chr.IP, chr.Port)

	client, err := cli.NewClient(ctx, chr.Addr(), logger)
	if err != nil {
		fatalf("could not get a client: %v", err)
	}

	fmt.Print("\nGetting receiver status...")
	status, err := command.Status.Get(client)
	if err != nil {
		fatalf("could not get status: %v", err)
	}
	fmt.Println(" OK")
	cli.FprintStatus(os.Stdout, status)

	// Get Media app
	fmt.Print("\nLooking for a media app...")
	app, err := media.FromStatus(client, status)
	if err == nil {
		fmt.Println(" OK")
		go app.UpdateStatus()
		st, err := app.GetStatus()
		if err != nil {
			fatalf("could not get media status: %v", err)
		}
		for _, s := range st {
			if s.Item != nil {
				fmt.Printf("  Item: %s\n", s.Item.ContentId)
				fmt.Printf("  Duration: %s\n", secondsToDuration(s.Item.Duration))
			}
			fmt.Printf("  Current Time: %s\n", secondsToDuration(s.CurrentTime))
			fmt.Printf("  State: %s\n", s.PlayerState)
		}
		session, err := app.CurrentSession()
		if err != nil {
			fatalf("could not get a session")
		}
		if true {
			ch, err := session.Seek(media.Seek(15 * time.Second))
			if err != nil {
				fatalf("could not pause")
			}
			fmt.Println(string(<-ch))
		}
	} else {
		fmt.Println(" not found")
	}
}
