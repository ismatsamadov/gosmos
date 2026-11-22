package main

import (
	"sync"
	"time"
)

// Bucket is one customer's token bucket.
type Bucket struct {
	Tokens float64   // tokens left, at most 10
	Last   time.Time // when Tokens was last brought up to date
}

// Limiter keeps a bucket per customer; it is safe for many requests at once.
type Limiter struct {
	now     func() time.Time
	mu      sync.Mutex
	buckets map[string]*Bucket
}
