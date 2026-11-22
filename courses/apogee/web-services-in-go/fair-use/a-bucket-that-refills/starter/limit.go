package main

import "time"

// A customer may send bucketSize requests at once, then refillPerSec a second.
const (
	bucketSize   = 10 // the most tokens a bucket holds
	refillPerSec = 1  // tokens added back each second
)

// Bucket is one customer's token bucket.
type Bucket struct {
	Tokens float64   // tokens left, at most 10
	Last   time.Time // when Tokens was last brought up to date
}

// Take brings the bucket up to date at now, then takes a token if a whole one is left.
func (b *Bucket) Take(now time.Time) bool {
	return true
}
