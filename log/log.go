package log

import (
	"io"

	kitlog "github.com/go-kit/kit/log"

	"github.com/oliverpool/go-chromecast"
)

// New creates a new structured logger.
func New(out io.Writer) chromecast.Logger {
	w := kitlog.NewSyncWriter(out)
	logger := kitlog.NewLogfmtLogger(w)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)
	return logger
}

// NopLogger returns a logger that doesn't do anything
func NopLogger() chromecast.Logger {
	return kitlog.NewNopLogger()
}
