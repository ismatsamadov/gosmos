package main

import (
	"math/rand/v2"
	"time"
)

// redialDelay is how long to wait after the attempt-th failed dial in a row, counted
// from 1: a second, doubled for each earlier failure, at most 30 seconds, times a jitter
// drawn from [0.5, 1.0).
func redialDelay(attempt int, r *rand.Rand) time.Duration {
	d := time.Second
	for i := 1; i < attempt && d < 30*time.Second; i++ {
		d *= 2
	}
	d = min(d, 30*time.Second)
	return time.Duration(float64(d) * (0.5 + r.Float64()/2))
}
