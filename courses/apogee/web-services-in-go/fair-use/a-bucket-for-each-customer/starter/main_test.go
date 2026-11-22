package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// clock is a fake clock the tests move by hand.
type clock struct {
	mu  sync.Mutex
	now time.Time
}

// Now returns the clock's time.
func (c *clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Add moves the clock forward by d.
func (c *clock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// newClock returns a clock that reads 09:00 UTC on Monday 10 January.
func newClock() *clock {
	return &clock{now: time.Date(2028, 1, 10, 9, 0, 0, 0, time.UTC)}
}

// allow calls l.Allow(customer) and writes down what it returned.
func allow(l *Limiter, customer string) string {
	ok, remaining, retryAfter := l.Allow(customer)
	return fmt.Sprintf("%v, %d, %v", ok, remaining, retryAfter)
}

func TestAllowCountsDown(t *testing.T) {
	l := newLimiter(newClock().Now)
	for i := 1; i <= 10; i++ {
		want := fmt.Sprintf("true, %d, 0s", 10-i)
		if i == 10 {
			want = "true, 0, 1s"
		}
		if got := allow(l, "brinewatch"); got != want {
			t.Errorf(`Allow("brinewatch") call %d = %s, want %s`, i, got, want)
		}
	}
	if got, want := allow(l, "brinewatch"), "false, 0, 1s"; got != want {
		t.Errorf(`Allow("brinewatch") call 11 = %s, want %s`, got, want)
	}
}

func TestAllowKeepsABucketPerCustomer(t *testing.T) {
	l := newLimiter(newClock().Now)
	for range 11 {
		l.Allow("brinewatch")
	}
	if got, want := allow(l, "cirrine"), "true, 9, 0s"; got != want {
		t.Errorf(`Allow("cirrine") after Brinewatch's 11 = %s, want %s`, got, want)
	}
	if got, want := allow(l, "brinewatch"), "false, 0, 1s"; got != want {
		t.Errorf(`Allow("brinewatch") after Cirrine's first = %s, want %s`, got, want)
	}
}

func TestAllowRefillsWithTheClock(t *testing.T) {
	c := newClock()
	l := newLimiter(c.Now)
	for range 11 {
		l.Allow("brinewatch")
	}
	steps := []struct {
		after time.Duration
		want  string
	}{
		{600 * time.Millisecond, "false, 0, 1s"}, // 0.6 of a token: 0.4 s to wait, rounded up
		{400 * time.Millisecond, "true, 0, 1s"},
		{2500 * time.Millisecond, "true, 1, 0s"},
		{time.Hour, "true, 9, 0s"},
	}
	for _, step := range steps {
		c.Add(step.after)
		if got := allow(l, "brinewatch"); got != step.want {
			t.Errorf(`Allow("brinewatch") at %s = %s, want %s`, c.Now().Format("15:04:05.000"), got, step.want)
		}
	}
}

func TestAllowManyAtOnce(t *testing.T) {
	l := newLimiter(newClock().Now)
	var allowed atomic.Int64
	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			for range 10 {
				if ok, _, _ := l.Allow("brinewatch"); ok {
					allowed.Add(1)
				}
			}
		})
	}
	wg.Wait()
	if got := allowed.Load(); got != 10 {
		t.Errorf(`Allow("brinewatch") from 20 goroutines, 10 calls each, at one instant = %d allowed, want 10`, got)
	}
}
