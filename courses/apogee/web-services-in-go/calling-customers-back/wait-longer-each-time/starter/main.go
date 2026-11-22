package main

import (
	"math/rand/v2"
	"time"
)

// backoff returns the wait after failed attempt number attempt: from 2s, doubling up to 5 minutes, with jitter from r.
func backoff(attempt int, r *rand.Rand) time.Duration {
	return 0
}
