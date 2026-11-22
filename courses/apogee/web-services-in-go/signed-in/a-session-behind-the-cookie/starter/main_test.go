package main

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

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

// Add moves the clock forward by d.
func (c *clock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// afternoon is when the tests start: Friday 21 January, 15:30 UTC.
var afternoon = time.Date(2028, 1, 21, 15, 30, 0, 0, time.UTC)

// random reports whether id looks like rand.Text's output: at least 26 characters, A to Z and 2 to 7.
func random(id string) bool {
	return len(id) >= 26 && strings.Trim(id, "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567") == ""
}

func TestStartGivesARandomID(t *testing.T) {
	s := newSessions((&clock{now: afternoon}).Now)
	first := s.Start("furrowlight")
	if !random(first) {
		t.Fatalf("Start(%q) = %q, want at least 26 random characters from A to Z and 2 to 7", "furrowlight", first)
	}
	if second := s.Start("furrowlight"); second == first {
		t.Errorf("two Starts for Furrowlight both gave %q, want two different IDs", first)
	}
}

func TestHollisEditsHisCookie(t *testing.T) {
	s := newSessions((&clock{now: afternoon}).Now)
	s.Start("furrowlight")
	for _, typed := range []string{"cirrine", "furrowlight", ""} {
		if customer, ok := s.Lookup(typed); ok {
			t.Errorf("Lookup(%q) = %q, %v, want \"\", false: a cookie anyone can type signs nobody in", typed, customer, ok)
		}
	}
}

func TestLookupFindsTheCustomer(t *testing.T) {
	s := newSessions((&clock{now: afternoon}).Now)
	hollis := s.Start("furrowlight")
	yusra := s.Start("apogee")
	if customer, ok := s.Lookup(hollis); customer != "furrowlight" || !ok {
		t.Errorf("Lookup(Hollis's ID) = %q, %v, want %q, true", customer, ok, "furrowlight")
	}
	if customer, ok := s.Lookup(yusra); customer != "apogee" || !ok {
		t.Errorf("Lookup(Yusra's ID) = %q, %v, want %q, true", customer, ok, "apogee")
	}
}

func TestUnusedSessionEnds(t *testing.T) {
	clk := &clock{now: afternoon}
	s := newSessions(clk.Now)
	id := s.Start("furrowlight")
	clk.Add(29*time.Minute + 59*time.Second)
	if _, ok := s.Lookup(id); !ok {
		t.Fatalf("Lookup(id) 29m59s after Start = false, want true: the session still lasts")
	}

	clk = &clock{now: afternoon}
	s = newSessions(clk.Now)
	id = s.Start("furrowlight")
	clk.Add(30 * time.Minute)
	if customer, ok := s.Lookup(id); ok {
		t.Errorf("Lookup(id) exactly 30m after Start = %q, %v, want \"\", false: the session has ended", customer, ok)
	}
	if _, kept := s.byID[id]; kept {
		t.Errorf("after Lookup found it ended, the session is still kept, want it deleted")
	}
}

func TestUseRenews(t *testing.T) {
	clk := &clock{now: afternoon}
	s := newSessions(clk.Now)
	id := s.Start("apogee")
	for _, step := range []time.Duration{20 * time.Minute, 25 * time.Minute} {
		clk.Add(step)
		if _, ok := s.Lookup(id); !ok {
			t.Fatalf("at %s, Lookup(id) = false, want true: it was last used %v before, and each use moves its end 30 minutes on",
				clk.Now().Format("15:04"), step)
		}
	}
	clk.Add(30 * time.Minute)
	if _, ok := s.Lookup(id); ok {
		t.Errorf("at %s, Lookup(id) = true, want false: its last use was exactly 30 minutes before", clk.Now().Format("15:04"))
	}
}

func TestManyBrowsersAtOnce(t *testing.T) {
	s := newSessions((&clock{now: afternoon}).Now)
	ids := make([]string, 50)
	var wg sync.WaitGroup
	for i := range ids {
		wg.Go(func() {
			ids[i] = s.Start(fmt.Sprintf("customer-%02d", i))
			s.Lookup(ids[i])
		})
	}
	wg.Wait()
	seen := map[string]bool{}
	for i, id := range ids {
		want := fmt.Sprintf("customer-%02d", i)
		if customer, ok := s.Lookup(id); customer != want || !ok {
			t.Errorf("after 50 Starts at once, Lookup(%q) = %q, %v, want %q, true", id, customer, ok, want)
		}
		if seen[id] {
			t.Errorf("after 50 Starts at once, %q was given twice", id)
		}
		seen[id] = true
	}
}
