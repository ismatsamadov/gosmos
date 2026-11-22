package main

import (
	"errors"
	"net/http"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	feb    = time.Date(2028, 2, 1, 0, 0, 0, 0, time.UTC)
	mar    = time.Date(2028, 3, 1, 0, 0, 0, 0, time.UTC)
	friday = time.Date(2028, 2, 4, 11, 0, 0, 0, time.UTC) // Yusra's TL;DR
)

// refusal is the message of a refused booking under a plan of plan bookings, whose count starts again at next.
func refusal(plan, next string) string {
	return "plan of " + plan + " bookings a month used up: the count starts again at " + next
}

func TestTake(t *testing.T) {
	dec := time.Date(2028, 12, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name     string
		u        usage
		customer string
		plan     int
		now      time.Time
		refused  string    // the refusal's message, "" when the booking counts
		month    time.Time // the month counted afterwards
		used     int       // the customer's count afterwards
	}{
		{"Brinewatch's 12th February booking", usage{feb, map[string]int{"brinewatch": 11}}, "brinewatch", 12, friday, "", feb, 12},
		{"its 13th", usage{feb, map[string]int{"brinewatch": 12}}, "brinewatch", 12, friday, refusal("12", "2028-03-01T00:00:00Z"), feb, 12},
		{"the last second of February", usage{feb, map[string]int{"brinewatch": 12}}, "brinewatch", 12, time.Date(2028, 2, 29, 23, 59, 59, 0, time.UTC), refusal("12", "2028-03-01T00:00:00Z"), feb, 12},
		{"the first second of March", usage{feb, map[string]int{"brinewatch": 12}}, "brinewatch", 12, mar, "", mar, 1},
		{"Cirrine's 41st", usage{feb, map[string]int{"cirrine": 40}}, "cirrine", 40, friday, refusal("40", "2028-03-01T00:00:00Z"), feb, 40},
		{"Cirrine beside Brinewatch's full plan", usage{feb, map[string]int{"brinewatch": 12, "cirrine": 4}}, "cirrine", 40, friday, "", feb, 5},
		{"a clock reading 21:30 at UTC-3 on 29 February", usage{feb, map[string]int{"brinewatch": 12}}, "brinewatch", 12, time.Date(2028, 2, 29, 21, 30, 0, 0, time.FixedZone("UTC-3", -3*60*60)), "", mar, 1},
		{"February a year on", usage{feb, map[string]int{"brinewatch": 12}}, "brinewatch", 12, time.Date(2029, 2, 10, 9, 0, 0, 0, time.UTC), "", time.Date(2029, 2, 1, 0, 0, 0, 0, time.UTC), 1},
		{"December's plan", usage{dec, map[string]int{"brinewatch": 12}}, "brinewatch", 12, time.Date(2028, 12, 31, 23, 0, 0, 0, time.UTC), refusal("12", "2029-01-01T00:00:00Z"), dec, 12},
		{"a count that never started", usage{}, "brinewatch", 12, friday, "", feb, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u := c.u
			err := u.take(c.customer, c.plan, c.now)
			call := "take(" + c.customer + ", at " + c.now.Format(time.RFC3339) + ")"
			if c.refused == "" && err != nil {
				t.Errorf("%s = %v, want nil", call, err)
			}
			if c.refused != "" {
				apiErr, ok := errors.AsType[*APIError](err)
				if !ok || apiErr.Status != http.StatusForbidden || apiErr.Code != "quota_exceeded" || apiErr.Message != c.refused {
					t.Errorf("%s = %#v, want a 403 quota_exceeded %q", call, err, c.refused)
				}
			}
			if !u.Month.Equal(c.month) || u.Month.Location() != time.UTC {
				t.Errorf("after %s, Month = %v, want %v", call, u.Month, c.month)
			}
			if got := u.Used[c.customer]; got != c.used {
				t.Errorf("after %s, Used[%q] = %d, want %d", call, c.customer, got, c.used)
			}
		})
	}
}

// clock is a fake clock that moves only when the test moves it.
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

// Set moves the clock to t.
func (c *clock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t
}

// plans are the customers' plans, in bookings a month. Apogee's own have none.
var plans = map[string]int{"furrowlight": 60, "cirrine": 40, "brinewatch": 12}

func TestTheCountSurvivesARestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "usage.json")
	clk := &clock{now: friday}
	q, err := openQuota(path, plans, clk.Now)
	if err != nil {
		t.Fatalf("openQuota = %v, want nil", err)
	}
	for i := 1; i <= 12; i++ {
		if err := q.Take("brinewatch"); err != nil {
			t.Fatalf("Brinewatch's booking %d of 12 = %v, want nil", i, err)
		}
	}
	if err := q.Take("brinewatch"); err == nil {
		t.Errorf("Brinewatch's 13th February booking = nil, want the refusal")
	}
	for i := 0; i < 100; i++ {
		if err := q.Take("apogee"); err != nil {
			t.Fatalf("Apogee's own booking %d = %v, want nil: Apogee's tests have no plan", i+1, err)
		}
	}
	q, err = openQuota(path, plans, clk.Now) // Friday's deploy
	if err != nil {
		t.Fatalf("openQuota after a restart = %v, want nil", err)
	}
	if err := q.Take("brinewatch"); err == nil {
		t.Errorf("Brinewatch's 13th February booking after a restart = nil, want the refusal: the count is in usage.json")
	}
	clk.Set(mar)
	if err := q.Take("brinewatch"); err != nil {
		t.Errorf("Brinewatch's first March booking = %v, want nil", err)
	}
}

func TestFiftyAtOnce(t *testing.T) {
	q, err := openQuota(filepath.Join(t.TempDir(), "usage.json"), plans, (&clock{now: friday}).Now)
	if err != nil {
		t.Fatalf("openQuota = %v, want nil", err)
	}
	for i := 0; i < 30; i++ {
		if err := q.Take("cirrine"); err != nil {
			t.Fatalf("Cirrine's booking %d = %v, want nil", i+1, err)
		}
	}
	var counted atomic.Int32
	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() {
			if q.Take("cirrine") == nil {
				counted.Add(1)
			}
		})
	}
	wg.Wait()
	if n := counted.Load(); n != 10 {
		t.Errorf("50 bookings at once with 30 of Cirrine's 40 made: %d counted, want 10", n)
	}
}
