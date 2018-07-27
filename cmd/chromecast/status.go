package main

import (
	"fmt"

	"github.com/oliverpool/go-chromecast/command/media"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the status of the first chromecast found",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, ctx, cancel := flags()
		defer cancel()

		client, status, err := GetClientWithStatus(ctx, logger)
		if err != nil {
			return fmt.Errorf("could not get a client: %v", err)
		}
		defer client.Close()
		fmt.Println("\n", status.String())

		// Get media app
		fmt.Print("\nLooking for a media app...")
		app, err := media.ConnectFromStatus(client, status)
		if err != nil {
			return fmt.Errorf("no media app found: %v", err)
		}
		fmt.Println(" OK")

		// Get media app status
		fmt.Print("Getting media app status...")
		st, err := app.Status()
		if err != nil {
			return fmt.Errorf("could not get media status: %v", err)
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
		return nil
	},
}
