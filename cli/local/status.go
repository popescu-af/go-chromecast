package local

import (
	"fmt"
	"sync"
	"time"

	chromecast "github.com/popescu-af/go-chromecast"
	"github.com/popescu-af/go-chromecast/command/media"
)

type Status struct {
	mu          sync.Mutex
	volume      float64
	muted       bool
	playerState string
	time        time.Duration
	totalTime   time.Duration
	orderSent   time.Time
}

func New(cstatus chromecast.Status) *Status {
	s := Status{}

	vol := cstatus.Volume
	if vol != nil {
		if vol.Level != nil {
			s.volume = *vol.Level
		}
		if vol.Muted != nil {
			s.muted = *vol.Muted
		}
	}
	return &s
}

func (s *Status) UpdateMedia(mstatus media.Status) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if time.Since(s.orderSent) < time.Second {
		return int(s.time.Seconds())
	}
	s.playerState = mstatus.PlayerState
	s.time = mstatus.CurrentTime.Duration
	if mstatus.Item != nil {
		s.totalTime = mstatus.Item.Duration.Duration
	}
	return int(mstatus.CurrentTime.Seconds())
}

func (s *Status) order() func() {
	s.mu.Lock()
	s.orderSent = time.Now()
	return s.mu.Unlock
}

func (s *Status) TogglePlay() bool {
	defer s.order()()

	if s.playerState == "PAUSED" {
		s.playerState = "PLAYING"
		return true
	}
	s.playerState = "PAUSED"
	return false
}

func (s *Status) ToggleMute() bool {
	defer s.order()()

	s.muted = !s.muted
	return s.muted
}

func (s *Status) IncrVolume(diff float64) float64 {
	defer s.order()()

	s.volume += diff
	if s.volume > 1 {
		s.volume = 1
	} else if s.volume < 0 {
		s.volume = 0
	}
	return s.volume
}

func (s *Status) SeekBy(diff time.Duration) time.Duration {
	defer s.order()()

	s.time += diff
	if s.time < 0 {
		s.time = 0
	}
	return s.time
}

func (s *Status) PlayerState() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.playerState {
	case "PLAYING":
		return " Playing "
	case "PAUSED":
		return "[paused] "
	}
	return s.playerState
}

func (s *Status) TimeStatus() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return fmt.Sprintf("%-8s/%8s", s.time.Round(time.Second), s.totalTime.Round(time.Second))
}
