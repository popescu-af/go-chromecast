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
	client, err := cli.GetClient(ctx, "", 0, "", logger)
	if err != nil {
		fatalf("could not get a client: %v", err)
	}
	fmt.Println(" OK")

	launcher := receiver.Launcher{
		Requester: client,
	}

	fmt.Print("\nGetting receiver status...")
	status, err := launcher.Status()
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
		st, err := app.Status()
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
