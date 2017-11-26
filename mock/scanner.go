package mock

import (
	cast "github.com/oliverpool/go-chromecast"
	"golang.org/x/net/context"
)

type Scanner struct {
	ScanFuncCalled int
	ScanFunc       func(ctx context.Context, results chan<- *cast.Device) error
}

func (s *Scanner) Scan(ctx context.Context, results chan<- *cast.Device) error {
	s.ScanFuncCalled++
	return s.ScanFunc(ctx, results)
}
