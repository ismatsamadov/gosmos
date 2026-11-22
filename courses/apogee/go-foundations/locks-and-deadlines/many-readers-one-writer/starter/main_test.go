package main

import (
	"fmt"
	"maps"
	"sync"
	"testing"
	"time"
)

// whileReading calls f while another reader holds c's lock, and stops the test if f waits for it.
func whileReading(t *testing.T, c *StatusCache, call string, f func()) {
	t.Helper()
	c.mu.RLock() // another service is reading the cache right now
	defer c.mu.RUnlock()
	done := make(chan struct{})
	go func() {
		f()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("%s waited 1s for another reader to finish: readers share the lock with RLock", call)
	}
}

func TestSnapshotIsACopy(t *testing.T) {
	c := NewStatusCache()
	c.Set("CWR", "tracking pipit-3")
	snap := c.Snapshot()
	snap["CWR"] = "stowed"
	if got := c.Get("CWR"); got != "tracking pipit-3" {
		t.Errorf(`after a change to a snapshot, Get("CWR") = %q, want "tracking pipit-3"`, got)
	}
}

func TestReadersShareTheLock(t *testing.T) {
	c := NewStatusCache()
	c.Set("YST", "stowed")
	whileReading(t, c, `Get("YST")`, func() { c.Get("YST") })
	whileReading(t, c, "Snapshot()", func() { c.Snapshot() })
}

func TestReadersWaitForTheWriter(t *testing.T) {
	reads := []struct {
		call string
		read func(c *StatusCache) string
	}{
		{`Get("CWR")`, func(c *StatusCache) string { return c.Get("CWR") }},
		{`Snapshot()["CWR"]`, func(c *StatusCache) string { return c.Snapshot()["CWR"] }},
	}
	for _, r := range reads {
		c := NewStatusCache()
		c.Set("CWR", "stowed")
		c.mu.Lock() // the poller starts writing
		got := make(chan string, 1)
		go func() {
			got <- r.read(c)
		}()
		time.Sleep(20 * time.Millisecond) // the read arrives in the middle of the write
		c.status["CWR"] = "tracking pipit-3"
		c.mu.Unlock() // the write is done
		select {
		case status := <-got:
			if status != "tracking pipit-3" {
				t.Errorf(`%s during a write = %q, want "tracking pipit-3": a reader waits for the writer's lock`, r.call, status)
			}
		case <-time.After(time.Second):
			t.Fatalf("%s has not returned 1s after the write ended", r.call)
		}
	}
}

func TestSetsAndSnapshotsTogether(t *testing.T) {
	c := NewStatusCache()
	codes := []string{"YST", "PEL", "SQT", "LAN", "CWR"}
	var wg sync.WaitGroup
	for _, code := range codes {
		wg.Go(func() { // the poller, writing one station's updates
			for i := range 100 {
				c.Set(code, fmt.Sprintf("tracking, frame %d", i))
			}
		})
		wg.Go(func() { // a service, reading
			for range 100 {
				c.Snapshot()
			}
		})
	}
	wg.Wait()
	want := map[string]string{}
	for _, code := range codes {
		want[code] = "tracking, frame 99"
	}
	if got := c.Snapshot(); !maps.Equal(got, want) {
		t.Errorf("after concurrent Sets and Snapshots, Snapshot() = %v, want %v", got, want)
	}
}
