package main

import (
	"sync"
	"testing"
	"time"
)

// within runs f, which calls Add, and stops the test if f hasn't returned within 1 s.
func within(t *testing.T, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("%s has not returned after 1s: does Add unlock c.mu before it returns?", what)
	}
}

func TestAddFromFiftyFeeds(t *testing.T) {
	c := NewCounter()
	sats := []string{"pipit-1", "pipit-2", "furrow-a", "brine-3", "cirrine-1"}
	within(t, "50 feeds of 200 Adds each", func() {
		var wg sync.WaitGroup
		for feed := range 50 {
			wg.Go(func() {
				for i := range 200 {
					c.Add(sats[(feed+i)%len(sats)])
				}
			})
		}
		wg.Wait()
	})
	for _, sat := range sats {
		if got := c.Get(sat); got != 2000 {
			t.Errorf("after 50 feeds of 200 Adds each, Get(%q) = %d, want 2000", sat, got)
		}
	}
}

func TestGetUnknownSatellite(t *testing.T) {
	c := NewCounter()
	var got int
	within(t, `Add("pipit-1") then Get("brine-2")`, func() {
		c.Add("pipit-1")
		got = c.Get("brine-2")
	})
	if got != 0 {
		t.Errorf(`Get("brine-2") = %d, want 0`, got)
	}
}
