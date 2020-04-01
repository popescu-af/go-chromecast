package main

import (
	"fmt"

	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/discovery/zeroconf"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all the chromecast found in the network",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, ctx, cancel := flags()
		defer cancel()

		scanner := zeroconf.Scanner{Logger: logger}
		devices := make(chan *chromecast.Device, 5)
		seen := make(map[string]bool)

		if err := scanner.Scan(ctx, devices); err != nil {
			return fmt.Errorf("could not start scanner: %w", err)
		}

		for d := range devices {
			if seen[d.ID()] {
				continue
			}
			seen[d.ID()] = true
			fmt.Printf("- %s [addr=%s; uuid=%s; type=%s; status=%s]\n",
				d.Name(), d.Addr(), d.ID(), d.Type(), d.Status())
		}
		return nil
	},
}
