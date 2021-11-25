package streak_test

import (
	"testing"
	"time"

	"github.com/popescu-af/go-chromecast/streak"
)

const baseTime = 10 * time.Millisecond
const timeMargin = 1 * time.Millisecond

func TestEmptyStreaker(t *testing.T) {
	s := streak.New(time.Second)
	for i := 0; i < 10; i++ {
		if f := s.UpdatedFactor(); f != 1 {
			t.Errorf("updated factor should always be 1, got %d", f)
		}
	}
}

func TestOneLevelStreaker(t *testing.T) {
	d := baseTime
	s := streak.New(time.Second, streak.Factor{
		After: d,
		Value: 2,
	})
	start := time.Now()
	for time.Now().Sub(start) <= d {
		if f := s.UpdatedFactor(); f != 1 {
			t.Fatalf("updated factor should be 1 at the beginning, got %d after %v", f, time.Now().Sub(start))
		}
	}
	time.Sleep(timeMargin)

	for i := 0; i < 10; i++ {
		if f := s.UpdatedFactor(); f != 2 {
			t.Fatalf("updated factor should be 2 afterwards, got %d after %v", f, time.Now().Sub(start))
		}
	}
}

func TestOneLevelReleasedStreaker(t *testing.T) {
	d := baseTime
	s := streak.New(baseTime/2, streak.Factor{
		After: d,
		Value: 2,
	})
	start := time.Now()
	for time.Now().Sub(start) <= d {
		if f := s.UpdatedFactor(); f != 1 {
			t.Fatalf("updated factor should be 1 at the beginning, got %d after %v", f, time.Now().Sub(start))
		}
	}
	time.Sleep(timeMargin)

	if f := s.UpdatedFactor(); f != 2 {
		t.Fatalf("updated factor should be 2 afterwards, got %d after %v", f, time.Now().Sub(start))
	}
	time.Sleep(baseTime)

	for i := 0; i < 10; i++ {
		if f := s.UpdatedFactor(); f != 1 {
			t.Fatalf("updated factor should be 1 again at the end, got %d", f)
		}
	}
}

func TestOneLevelReleasedImmediateStreaker(t *testing.T) {
	d := baseTime
	s := streak.New(baseTime/2, streak.Factor{
		After: d,
		Value: 2,
	})
	start := time.Now()
	for time.Now().Sub(start) <= d {
		if f := s.UpdatedFactor(); f != 1 {
			t.Fatalf("updated factor should be 1 at the beginning, got %d after %v", f, time.Now().Sub(start))
		}
	}
	time.Sleep(timeMargin + baseTime)

	for i := 0; i < 10; i++ {
		if f := s.UpdatedFactor(); f != 1 {
			t.Fatalf("updated factor should be 1 again at the end, got %d", f)
		}
	}
}

func TestTwoLevelsReleasedStreaker(t *testing.T) {
	d := baseTime
	s := streak.New(baseTime/2, streak.Factor{
		After: d,
		Value: 2,
	}, streak.Factor{
		After: 2 * d,
		Value: 3,
	})
	start := time.Now()
	for time.Now().Sub(start) <= d {
		if f := s.UpdatedFactor(); f != 1 {
			t.Fatalf("updated factor should be 1 at the beginning, got %d after %v", f, time.Now().Sub(start))
		}
	}
	time.Sleep(timeMargin)

	for time.Now().Sub(start) <= 2*d {
		if f := s.UpdatedFactor(); f != 2 {
			t.Fatalf("updated factor should be 2 at level 2, got %d after %v", f, time.Now().Sub(start))
		}
	}
	time.Sleep(timeMargin)

	if f := s.UpdatedFactor(); f != 3 {
		t.Fatalf("updated factor should be 3 afterwards, got %d after %v", f, time.Now().Sub(start))
	}
	time.Sleep(baseTime)

	for i := 0; i < 10; i++ {
		if f := s.UpdatedFactor(); f != 1 {
			t.Fatalf("updated factor should be 1 again at the end, got %d", f)
		}
	}
}

func TestTwoLevelsReleasedStreakerReversed(t *testing.T) {
	d := baseTime
	s := streak.New(baseTime/2, streak.Factor{
		After: 2 * d,
		Value: 3,
	}, streak.Factor{
		After: d,
		Value: 2,
	})
	start := time.Now()
	for time.Now().Sub(start) <= d {
		if f := s.UpdatedFactor(); f != 1 {
			t.Fatalf("updated factor should be 1 at the beginning, got %d after %v", f, time.Now().Sub(start))
		}
	}
	time.Sleep(timeMargin)

	for time.Now().Sub(start) <= 2*d {
		if f := s.UpdatedFactor(); f != 2 {
			t.Fatalf("updated factor should be 2 at level 2, got %d after %v", f, time.Now().Sub(start))
		}
	}
	time.Sleep(timeMargin)

	if f := s.UpdatedFactor(); f != 3 {
		t.Fatalf("updated factor should be 3 afterwards, got %d after %v", f, time.Now().Sub(start))
	}
	time.Sleep(baseTime)

	for i := 0; i < 10; i++ {
		if f := s.UpdatedFactor(); f != 1 {
			t.Fatalf("updated factor should be 1 again at the end, got %d", f)
		}
	}
}
