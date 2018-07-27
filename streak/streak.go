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

// Streaker contains the streaker parameters
// Factors must be sorted
type Streaker struct {
	MaxDelayBetweenHits time.Duration
	Factors             []Factor

	previousHit time.Time
	streakStart time.Time
}

// Factor indicates which factor should be returned after a streak lasted the given duration
type Factor struct {
	After time.Duration
	Value int64
}

// UpdatedFactor returns the factor after considering a hit
func (s *Streaker) UpdatedFactor() (value int64) {
	return s.Factor(s.Hit())
}

// Factor returns the factor for a given streak duration
func (s Streaker) Factor(streakDuration time.Duration) (value int64) {
	value = 1
	for _, f := range s.Factors {
		if f.After > streakDuration {
			break
		}
		value = f.Value
	}
	return value
}

// Hit considers a hit
// it returns the current streak duration
func (s *Streaker) Hit() time.Duration {
	now := time.Now()
	if now.Sub(s.previousHit) > s.MaxDelayBetweenHits {
		s.streakStart = now
	}
	s.previousHit = now
	return now.Sub(s.streakStart)
}
