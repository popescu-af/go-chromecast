package streak

import (
	"time"
)

// New creates a new Streaker
// A streaker returns a given Factor.Value after the Factor.Delay has passed (since the start of the streak)
// The value is reset to 1 when two hits are more than maxDelayBetweenHits apart.
// The factors will be automatically sorted.
func New(maxDelayBetweenHits time.Duration, factors ...Factor) Streaker {
	sortFactors(factors)

	return Streaker{
		MaxDelayBetweenHits: maxDelayBetweenHits,
		Factors:             factors,
	}
}

type Streaker struct {
	MaxDelayBetweenHits time.Duration
	Factors             []Factor

	previousHit time.Time
	streakStart time.Time
}

type Factor struct {
	After time.Duration
	Value int64
}

func (s *Streaker) UpdatedFactor() (value int64) {
	now := time.Now()

	if now.Sub(s.previousHit) > s.MaxDelayBetweenHits {
		s.streakStart = now
	}

	streakDuration := now.Sub(s.streakStart)
	value = 1
	for _, f := range s.Factors {
		if f.After > streakDuration {
			break
		}
		value = f.Value
	}

	s.previousHit = now
	return value
}
