package main

import (
	"context"
	"fmt"
	"os"

	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/discovery"
	"github.com/oliverpool/go-chromecast/discovery/zeroconf"
	"github.com/oliverpool/go-chromecast/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all the chromecast found in the network",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.New(os.Stdout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		scanner := zeroconf.Scanner{Logger: logger}
		allDevices := make(chan *chromecast.Device, 5)
		uniqDevices := make(chan *chromecast.Device, 5)

		worker, err := scanner.Scan(ctx, allDevices)
		if err != nil {
			return fmt.Errorf("could not initialize scanner: %v", err)
		}
		go worker()
		go discovery.Uniq(allDevices, uniqDevices)
		for d := range uniqDevices {
			fmt.Printf("- %s [Addr=%s; uuid=%s; type=%s; status=%s]\n",
				d.Name(), d.Addr(), d.ID(), d.Type(), d.Status())
		}
		return nil
	},
}
