package main

import (
	"fmt"
	"os"
	"time"

	"context"

	"github.com/codegangsta/cli"
	"github.com/oliverpool/go-chromecast"
	clicast "github.com/oliverpool/go-chromecast/cli"
	"github.com/oliverpool/go-chromecast/client"
	"github.com/oliverpool/go-chromecast/command"
	"github.com/oliverpool/go-chromecast/log"
)

func checkErr(err error) {
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("Timeout exceeded")
		} else {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}

func main() {
	commonFlags := []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "enable debug logging",
		},
		cli.StringFlag{
			Name:  "host",
			Usage: "chromecast hostname or IP",
		},
		cli.IntFlag{
			Name:  "port",
			Usage: "chromecast port",
			Value: 8009,
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "chromecast name",
		},
		cli.DurationFlag{
			Name:  "timeout",
			Value: 15 * time.Second,
		},
	}
	app := cli.NewApp()
	app.Name = "cast"
	app.Usage = "Command line tool for the Chromecast"
	app.Version = cast.Version
	app.Flags = commonFlags
	app.Commands = []cli.Command{
		{
			Name:   "status",
			Usage:  "Get status of the Chromecast",
			Action: statusCommand,
		},
		{
			Name:   "discover",
			Usage:  "Discover Chromecast devices",
			Action: discoverCommand,
		},
	}
	app.Run(os.Args)
	log.Println("Done")
}

// clientFromContext will try to get a cast client.
// If host is set, it will be used (along port).
// Otherwise, if name is set, a chromecast will be looked-up by name.
// Otherwise the first chromecast found will be returned.
func clientFromContext(ctx context.Context, c *cli.Context) *client.Client {
	chr, err := clicast.GetDevice(
		ctx,
		c.GlobalString("host"),
		c.GlobalInt("port"),
		c.GlobalString("name"),
	)
	checkErr(err)
	fmt.Printf("Found '%s' (%s:%d)...\n", chr.Name(), chr.IP, chr.Port)

	client, err := clicast.NewClient(ctx, chr.Addr())
	checkErr(err)
	return client
}

func statusCommand(c *cli.Context) {
	log.Debug = c.GlobalBool("debug")
	ctx, cancel := context.WithTimeout(context.Background(), c.GlobalDuration("timeout"))
	defer cancel()

	client := clientFromContext(ctx, c)

	// Get status
	fmt.Println("Status:")
	status, err := command.Status.Get(client.Request)
	checkErr(err)

	clicast.FprintStatus(os.Stdout, status)
}

func discoverCommand(c *cli.Context) {
	log.Debug = c.GlobalBool("debug")
	timeout := c.GlobalDuration("timeout")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Printf("Running scanner for %s...\n", timeout)
	for client := range clicast.Scan(ctx) {
		fmt.Printf("Found: %s:%d '%s' (%s: %s) %s\n", client.IP, client.Port, client.Name(), client.Type(), client.ID(), client.Status())
	}
	fmt.Println("Done")
}
