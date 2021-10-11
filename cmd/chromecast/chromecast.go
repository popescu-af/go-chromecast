package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oliverpool/go-chromecast"

	"github.com/oliverpool/go-chromecast/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chromecast",
	Short: "chromecast allows you to interact with a Chromecast",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("control")
			return controlCmd.RunE(cmd, args)
		}
		logger, ctx, cancel := flags(cmd)
		defer cancel()

		url := strings.TrimSpace(args[0])
		if _, err := os.Stat(url); err == nil {
			url, err = serveLocalFile(url)
			if err != nil {
				return fmt.Errorf("could not serve: %w", err)
			}
		}
		fmt.Println("Loading", url)

		return loadURL(ctx, cancel, url, "default", true, logger)
	},
}

func serveLocalFile(path string) (string, error) {
	ip, err := getOutboundIP()
	if err != nil {
		return "", err
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}
	s := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path)
		}),
	}
	go s.Serve(listener)

	filename := filepath.Base(path)
	port := listener.Addr().(*net.TCPAddr).Port
	addr := fmt.Sprintf("http://%s:%d/%s", ip, port, filename)
	return addr, nil
}

var timeout time.Duration
var verbose bool

func flags(cmd *cobra.Command) (chromecast.Logger, context.Context, context.CancelFunc) {
	cmd.Root().SilenceUsage = true
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
