package main

import (
	"net"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// Every test runs in a synctest bubble, whose clock is fake: the barrier's patience and every
// guard pass in no real time, and time moves on only once every goroutine is waiting, so a
// dish asked in a goroutine of its own always reaches the barrier in time.

func TestPollAll(t *testing.T) {
	t.Run("both dishes, asked together", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			b := newBarrier(2)
			dish1, dish2 := yst1(), yst2()
			dish1.Before, dish2.Before = b.hold(100*time.Millisecond), b.hold(0) // Dish 2 answers first
			clients := []*Client{connect(t, dish1), connect(t, dish2)}
			got, ok := pollWithin(clients, 3*time.Second)
			if !ok {
				t.Fatalf("pollAll(YST-1, YST-2) has not returned after 3s")
			}
			if b.waited() {
				t.Errorf("pollAll(YST-1, YST-2) asked one dish only after the other had answered, and the simulators gave up waiting after 500ms: poll each dish in its own goroutine")
			}
			want := []Result{{Name: "YST-1", Status: status1}, {Name: "YST-2", Status: status2}}
			if len(got) != len(want) {
				t.Fatalf("pollAll(YST-1, YST-2) returned %d results, want %d", len(got), len(want))
			}
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("pollAll(YST-1, YST-2)[%d] = %+v, want %+v: store each result at its client's index", i, got[i], want[i])
				}
			}
		})
	})
	t.Run("Dish 2 first", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			b := newBarrier(2)
			dish1, dish2 := yst1(), yst2()
			dish1.Before, dish2.Before = b.hold(0), b.hold(100*time.Millisecond) // Dish 1 answers first
			clients := []*Client{connect(t, dish2), connect(t, dish1)}
			got, ok := pollWithin(clients, 3*time.Second)
			if !ok || len(got) != 2 || got[0].Name != "YST-2" || got[1].Name != "YST-1" {
				t.Errorf("pollAll(YST-2, YST-1) = %+v, want YST-2's result, then YST-1's", got)
			}
		})
	})
	t.Run("one dish", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			got, ok := pollWithin([]*Client{connect(t, yst1())}, 3*time.Second)
			if !ok || len(got) != 1 || got[0] != (Result{Name: "YST-1", Status: status1}) {
				t.Errorf("pollAll(YST-1) = %+v, want [%+v]", got, Result{Name: "YST-1", Status: status1})
			}
		})
	})
	t.Run("no dishes", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			if got, ok := pollWithin(nil, 3*time.Second); !ok || len(got) != 0 {
				t.Errorf("pollAll(no clients) = %+v, want no results", got)
			}
		})
	})
}

// The two dishes on Monday 24 April: one anemometer, so one wind.
var (
	status1 = Status{Dish: 1, State: Tracking, AzDeg: 183.25, ElDeg: 41.70, WindKmh: 64, TempC: -12.5, Flags: 1, LastSeq: 118}
	status2 = Status{Dish: 2, State: Idle, AzDeg: 180, ElDeg: 45, WindKmh: 64, TempC: -12.5, Flags: 1}
)

// yst1 is a simulator for Dish 1, tracking at 183.25° and 41.70°.
func yst1() *Sim {
	sim := NewSim("YST-1", 1)
	sim.State, sim.AzDeg, sim.ElDeg, sim.Wind, sim.LastSeq = 2, 183.25, 41.70, 64, 118
	return sim
}

// yst2 is a simulator for Dish 2, idle at 180.00° and 45.00°.
func yst2() *Sim {
	sim := NewSim("YST-2", 2)
	sim.Wind = 64
	return sim
}

// barrier holds each simulator's answer to a QUERY until n dishes have been asked, as the
// backhaul's round trips would hold dishd if it asked one dish after the other. After 500ms
// it gives up waiting, and remembers that it did.
type barrier struct {
	n      int
	all    chan struct{}
	mu     sync.Mutex
	asked  int
	gaveUp bool
}

// newBarrier returns a barrier for n dishes.
func newBarrier(n int) *barrier {
	return &barrier{n: n, all: make(chan struct{})}
}

// hold returns a simulator's Before hook: it holds the answer to a QUERY at the barrier, then
// for extra more, so the dishes answer in the order the test chooses.
func (b *barrier) hold(extra time.Duration) func(req string) {
	return func(req string) {
		if req != "QUERY" {
			return
		}
		b.mu.Lock()
		b.asked++
		if b.asked == b.n {
			close(b.all)
		}
		b.mu.Unlock()
		select {
		case <-b.all:
		case <-time.After(500 * time.Millisecond):
			b.mu.Lock()
			b.gaveUp = true
			b.mu.Unlock()
		}
		time.Sleep(extra)
	}
}

// waited reports whether the barrier ever gave up waiting.
func (b *barrier) waited() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.gaveUp
}

// connect serves sim on one end of a pipe and returns a client on the other end that has said
// HELLO. When the test ends, both ends close and the simulator stops.
func connect(t *testing.T, sim *Sim) *Client {
	t.Helper()
	end, simEnd := net.Pipe()
	stopped := make(chan struct{})
	go func() {
		sim.ServeConn(simEnd)
		close(stopped)
	}()
	t.Cleanup(func() {
		end.Close()
		simEnd.Close()
		<-stopped
	})
	c, err := newClient(end)
	if err != nil {
		t.Fatalf("newClient(simulator) = %v, want a client", err)
	}
	if _, err := c.hello(); err != nil {
		t.Fatalf("hello() = %v, want a WELCOME", err)
	}
	return c
}

// pollWithin calls pollAll and waits for it at most limit, reporting whether it returned.
func pollWithin(clients []*Client, limit time.Duration) ([]Result, bool) {
	got := make(chan []Result, 1)
	go func() { got <- pollAll(clients) }()
	select {
	case r := <-got:
		return r, true
	case <-time.After(limit):
		for _, c := range clients {
			c.Close()
		}
		<-got
		return nil, false
	}
}
