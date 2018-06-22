package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/oliverpool/go-chromecast"

	"github.com/oliverpool/go-chromecast/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chromecast",
	Short: "chromecast allows you to interact with a Chromecast",
}

var timeout time.Duration
var verbose bool
var deviceName string
var deviceID string
var deviceIP net.IP
var devicePort int

func flags() (chromecast.Logger, context.Context, context.CancelFunc) {
	rootCmd.SilenceUsage = true
	logger := log.NopLogger()
	if verbose {
		logger = log.New(os.Stdout)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return logger, ctx, cancel
}

func init() {
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "Duration before stopping looking for chromecast(s)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose (debug) output")
	rootCmd.PersistentFlags().IPVar(&deviceIP, "ip", nil, "Specify chromecast IP")
	rootCmd.PersistentFlags().IntVar(&devicePort, "port", 8009, "Specify chromecast port (ignored if IP is not set)")
	rootCmd.PersistentFlags().StringVarP(&deviceName, "name", "n", "", "Specify chromecast name (ignored if IP is set)")
	rootCmd.PersistentFlags().StringVar(&deviceID, "id", "", "Specify chromecast ID (ignored if IP is set)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
