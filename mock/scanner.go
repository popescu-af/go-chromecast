package mock

import (
	"github.com/oliverpool/go-chromecast"
	"context"
)

type Scanner struct {
	ScanFuncCalled int
	ScanFunc       func(ctx context.Context, results chan<- *chromecast.Device) error
}

func (s *Scanner) Scan(ctx context.Context, results chan<- *chromecast.Device) error {
	s.ScanFuncCalled++
	return s.ScanFunc(ctx, results)
}
