package main

import (
	"context"
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

func flags() (chromecast.Logger, context.Context, context.CancelFunc) {
	rootCmd.SilenceUsage = true
	logger := log.NopLogger()
	if verbose {
		logger = log.New(os.Stdout)
	}
	ctx := context.Background()
	cancel := func() {}
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}
	return logger, ctx, cancel
}

func init() {
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "Duration before stopping looking for chromecast(s) (0 means no timeout)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose (debug) output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
