package main

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

// guard is how long a test waits for what a right stage does at once: one second in all.
type guard struct{ end time.Time }

// newGuard starts a guard's second.
func newGuard() guard { return guard{end: time.Now().Add(time.Second)} }

// wait is ready when the guard's second is over, and never sooner than 50 ms from now.
func (g guard) wait() <-chan time.Time {
	return time.After(max(time.Until(g.end), 50*time.Millisecond))
}

// boardQueue returns the Board's queue, with room for four schedules, holding the given versions.
func boardQueue(versions ...uint64) chan *Schedule {
	q := make(chan *Schedule, 4)
	for _, v := range versions {
		q <- &Schedule{Version: v}
	}
	return q
}

// waiting takes every schedule waiting in q, without waiting for more, and returns their versions in order.
func waiting(q chan *Schedule) []uint64 {
	var got []uint64
	for {
		select {
		case s := <-q:
			got = append(got, s.Version)
		default:
			return got
		}
	}
}

// offerAll offers the versions to q in turn, from a goroutine of its own, and returns how many offer dropped in all, or reports that an offer is still waiting once the guard ends.
func offerAll(g guard, q chan *Schedule, versions ...uint64) (dropped int, returned bool) {
	done := make(chan int, 1)
	go func() {
		n := 0
		for _, v := range versions {
			n += offer(q, &Schedule{Version: v})
		}
		done <- n
	}()
	select {
	case n := <-done:
		return n, true
	case <-g.wait():
		return 0, false
	}
}

func TestOffer(t *testing.T) {
	g := newGuard()
	tests := []struct {
		name    string
		queue   []uint64 // versions waiting before the offers
		offers  []uint64
		dropped int
		want    []uint64 // versions waiting after them
	}{
		{"an empty queue takes version 1", nil, []uint64{1}, 0, []uint64{1}},
		{"a full queue of versions 1 to 4, then version 5", []uint64{1, 2, 3, 4}, []uint64{5}, 1, []uint64{2, 3, 4, 5}},
		{"ten schedules while the Board's screen hangs", nil, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 6, []uint64{7, 8, 9, 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := boardQueue(tt.queue...)
			call := fmt.Sprintf("offer(q holding %v, versions %v)", tt.queue, tt.offers)
			dropped, returned := offerAll(g, q, tt.offers...)
			if !returned {
				t.Fatalf("%s is still waiting for the Board, want it to return at once", call)
			}
			if dropped != tt.dropped {
				t.Errorf("%s dropped %d, want %d", call, dropped, tt.dropped)
			}
			if got := waiting(q); !slices.Equal(got, tt.want) {
				t.Errorf("%s leaves versions %v waiting, want %v", call, got, tt.want)
			}
		})
	}
}
