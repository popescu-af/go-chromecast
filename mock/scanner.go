package mock

import (
	"context"

	"github.com/popescu-af/go-chromecast"
)

type Scanner struct {
	ScanFuncCalled int
	ScanFunc       func(ctx context.Context, results chan<- *chromecast.Device) (func() error, error)
}

func (s *Scanner) Scan(ctx context.Context, results chan<- *chromecast.Device) (func() error, error) {
	s.ScanFuncCalled++
	return s.ScanFunc(ctx, results)
}
