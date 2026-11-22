package main

import (
	"math/rand/v2"
	"time"
)

// backoff returns the wait after failed attempt number attempt: from 2s, doubling up to 5 minutes, with jitter from r.
func backoff(attempt int, r *rand.Rand) time.Duration {
	wait := 2 * time.Second
	for i := 1; i < attempt && wait < 5*time.Minute; i++ {
		wait *= 2
	}
	wait = min(wait, 5*time.Minute)
	return time.Duration(float64(wait) * (0.5 + r.Float64()/2))
}
