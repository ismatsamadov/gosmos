package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// sats are the seven satellites a station's uplinkd counts passes for.
var sats = []string{"pipit-1", "pipit-2", "pipit-3", "furrow-a", "furrow-b", "brine-3", "cirrine-1"}

// outcomes are the ways a pass can end.
var outcomes = []string{"flown", "failed", "skipped"}

// with calls v.With(a, b), and stops the test if it hasn't returned within a second, as when
// With asks for Lock while it still holds RLock.
func with(t *testing.T, v *CounterVec, a, b string) *Counter {
	t.Helper()
	got := make(chan *Counter, 1)
	go func() { got <- v.With(a, b) }()
	select {
	case c := <-got:
		return c
	case <-time.After(time.Second):
		t.Fatalf("With(%q, %q) still waiting after a second: does it take Lock while it holds RLock?", a, b)
		return nil
	}
}

// wait waits for wg, and stops the test if the goroutines haven't ended within a second.
func wait(t *testing.T, wg *sync.WaitGroup) {
	t.Helper()
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("goroutines calling With still waiting after a second: does it take Lock while it holds RLock?")
	}
}

// TestWithFromManyGoroutines has 40 goroutines count into every cell of the table at once,
// i+j+1 times each for satellite i and outcome j, so every cell has its own exact total.
func TestWithFromManyGoroutines(t *testing.T) {
	var v CounterVec
	var wg sync.WaitGroup
	for range 40 {
		wg.Go(func() {
			for i, sat := range sats {
				for j, outcome := range outcomes {
					for range i + j + 1 {
						v.With(sat, outcome).Inc()
					}
				}
			}
		})
	}
	wait(t, &wg)
	snap := v.Snapshot()
	var wrong []string
	for i, sat := range sats {
		for j, outcome := range outcomes {
			if got, want := snap[sat][outcome], uint64(40*(i+j+1)); got != want {
				wrong = append(wrong, fmt.Sprintf("Snapshot()[%q][%q] after 40 goroutines = %d, want %d", sat, outcome, got, want))
			}
		}
	}
	if len(wrong) > 0 {
		t.Errorf("%s (%d of 21 cells wrong)", wrong[0], len(wrong))
	}
}

func TestWithSamePair(t *testing.T) {
	var v CounterVec
	failed := with(t, &v, "brine-3", "failed")
	if again := with(t, &v, "brine-3", "failed"); again != failed {
		t.Errorf(`With("brine-3", "failed") twice = two different counters, want the same *Counter`)
	}
	if flown := with(t, &v, "brine-3", "flown"); flown == failed {
		t.Errorf(`With("brine-3", "flown") = the counter of "failed", want one of its own`)
	}
}

// TestWithOnePairAtOnce holds the map's lock while 40 goroutines ask for one new pair, then
// lets them in together: a With that doesn't look again under Lock adds a counter for each of
// them, and all but one lose their counts.
func TestWithOnePairAtOnce(t *testing.T) {
	var v CounterVec
	var asking sync.WaitGroup
	asking.Add(40)
	got := make([]*Counter, 40)
	var wg sync.WaitGroup
	v.mu.Lock()
	for i := range got {
		wg.Go(func() {
			asking.Done()
			got[i] = v.With("brine-3", "failed")
			got[i].Inc()
		})
	}
	asking.Wait()
	time.Sleep(20 * time.Millisecond) // time for each of them to reach RLock and wait there
	v.mu.Unlock()
	wait(t, &wg)
	distinct := map[*Counter]bool{}
	for _, c := range got {
		distinct[c] = true
	}
	if len(distinct) != 1 {
		t.Errorf(`With("brine-3", "failed") from 40 goroutines at once = %d different counters, want 1`, len(distinct))
	}
	if n := v.Snapshot()["brine-3"]["failed"]; n != 40 {
		t.Errorf(`Snapshot()["brine-3"]["failed"] after 40 Inc() = %d, want 40`, n)
	}
}

// TestSnapshotIsACopy counts Brine-3's passes, then changes a snapshot and counts more:
// neither may reach the other.
func TestSnapshotIsACopy(t *testing.T) {
	var v CounterVec
	with(t, &v, "brine-3", "flown").Add(23)
	with(t, &v, "brine-3", "failed").Inc()
	first := v.Snapshot()
	if got := first["brine-3"]["flown"]; got != 23 {
		t.Fatalf(`Snapshot()["brine-3"]["flown"] after Add(23) = %d, want 23`, got)
	}
	first["brine-3"]["flown"] = 0
	with(t, &v, "pipit-1", "flown").Inc()
	second := v.Snapshot()
	if got := second["brine-3"]["flown"]; got != 23 {
		t.Errorf(`Snapshot()["brine-3"]["flown"] after changing an earlier snapshot = %d, want 23`, got)
	}
	if got := second["brine-3"]["failed"]; got != 1 {
		t.Errorf(`Snapshot()["brine-3"]["failed"] after Inc() = %d, want 1`, got)
	}
	if _, ok := first["pipit-1"]; ok {
		t.Errorf(`an earlier Snapshot() gained "pipit-1" after a later With, want it unchanged`)
	}
	if n, ok := second["brine-3"]["skipped"]; ok {
		t.Errorf(`Snapshot()["brine-3"]["skipped"] = %d, want no entry for a pair never used`, n)
	}
}
