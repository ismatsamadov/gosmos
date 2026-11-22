package main

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"
)

func TestBackoffDoubles(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	same := rand.New(rand.NewPCG(1, 2)) // the test's own draws, from the same seed
	for attempt := 1; attempt <= 7; attempt++ {
		t.Run(fmt.Sprintf("after attempt %d", attempt), func(t *testing.T) {
			full := 2 * time.Second << (attempt - 1)
			want := time.Duration(float64(full) * (0.5 + same.Float64()/2))
			got := backoff(attempt, r)
			if got < full/2 || got >= full {
				t.Errorf("backoff(%d) = %v, want from %v up to %v", attempt, got, full/2, full)
			}
			if diff := got - want; diff < -time.Microsecond || diff > time.Microsecond {
				t.Errorf("backoff(%d) = %v, want %v: one r.Float64() for each wait, as 0.5 + r.Float64()/2", attempt, got, want)
			}
		})
	}
}

func TestBackoffStopsAtFiveMinutes(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	if got := backoff(8, r); got < 128*time.Second || got >= 256*time.Second {
		t.Errorf("backoff(8) = %v, want from 2m8s up to 4m16s", got)
	}
	for _, attempt := range []int{9, 20, 70} {
		if got := backoff(attempt, r); got < 150*time.Second || got >= 300*time.Second {
			t.Errorf("backoff(%d) = %v, want from 2m30s up to 5m0s", attempt, got)
		}
	}
}

func TestBackoffNeverZero(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for attempt := 1; attempt <= 100; attempt++ {
		if got := backoff(attempt, r); got <= 0 {
			t.Fatalf("backoff(%d) = %v, want a wait of at least a second", attempt, got)
		}
	}
}
