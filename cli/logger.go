package cli

import (
	"io"
	"os"

	kitlog "github.com/go-kit/kit/log"

	"github.com/oliverpool/go-chromecast"
)

// NewLogger creates a new structures logger.
func NewLogger(out io.Writer) chromecast.Logger {
	if os.Getenv("DEBUG") == "" {
		return kitlog.NewNopLogger()
	}
	w := kitlog.NewSyncWriter(out)
	logger := kitlog.NewLogfmtLogger(w)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)
	return logger
}
