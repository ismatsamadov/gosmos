package main

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

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
	if elapsed := now.Sub(b.Last); elapsed > 0 {
		b.Tokens = min(b.Tokens+elapsed.Seconds()*refillPerSec, bucketSize)
		b.Last = now
	}
	if b.Tokens < 1 {
		return false
	}
	b.Tokens--
	return true
}

// Limiter keeps a bucket per customer; it is safe for many requests at once.
type Limiter struct {
	now     func() time.Time
	mu      sync.Mutex
	buckets map[string]*Bucket
}

// newLimiter returns a Limiter with no buckets yet, which reads the time from now.
func newLimiter(now func() time.Time) *Limiter {
	return &Limiter{now: now, buckets: map[string]*Bucket{}}
}

// Allow takes a token from customer's bucket, and says how many whole tokens are left and, when none is, how long until one is back.
func (l *Limiter) Allow(customer string) (ok bool, remaining int, retryAfter time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	b := l.buckets[customer]
	if b == nil {
		b = &Bucket{}
		l.buckets[customer] = b
	}
	ok = b.Take(l.now())
	if b.Tokens < 1 {
		retryAfter = time.Duration(math.Ceil((1-b.Tokens)/refillPerSec)) * time.Second
	}
	return ok, int(b.Tokens), retryAfter
}

// errRateLimited is the answer to a request over its customer's limit.
var errRateLimited = &APIError{Status: http.StatusTooManyRequests, Code: "rate_limited", Message: "slow down: 10 requests at once, then 1 a second"}

// limit refuses a customer's request with 429 once its bucket is empty, and says on every answer how many are left.
func limit(l *Limiter) func(http.Handler) http.Handler {
	_ = errors.New
	_ = strconv.Itoa
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: answer 500 without a customer, and 429 with Retry-After once its bucket is empty.
			next.ServeHTTP(w, r)
		})
	}
}
