package main

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func TestPollAll(t *testing.T) {
	tests := []struct {
		name         string
		down1, down2 bool // the dish's controller hangs up when asked for its status
		want         []string
	}{
		{"both dishes answer", false, false, []string{
			"YST-1 TRACKING az=183.25 el=41.70 wind=64",
			"YST-2 IDLE az=180.00 el=45.00 wind=64",
		}},
		{"Dish 2's controller reboots", false, true, []string{
			"YST-1 TRACKING az=183.25 el=41.70 wind=64",
			"YST-2 DOWN closed",
		}},
		{"Dish 1's controller reboots", true, false, []string{
			"YST-1 DOWN closed",
			"YST-2 IDLE az=180.00 el=45.00 wind=64",
		}},
		{"both controllers reboot", true, true, []string{
			"YST-1 DOWN closed",
			"YST-2 DOWN closed",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clients := []*Client{connect(t, yst1(), tt.down1), connect(t, yst2(), tt.down2)}
			got, ok := pollWithin(clients, 3*time.Second)
			if !ok {
				t.Fatalf("pollAll(YST-1, YST-2) has not returned after 3s")
			}
			if len(got) != 2 {
				t.Fatalf("pollAll(YST-1, YST-2) returned %d results, want 2, one per dish, failed or not", len(got))
			}
			for i, down := range []bool{tt.down1, tt.down2} {
				want := Result{Name: []string{"YST-1", "YST-2"}[i], Status: []Status{status1, status2}[i]}
				if down {
					if got[i].Name != want.Name || !errors.Is(got[i].Err, io.EOF) {
						t.Errorf("pollAll(YST-1, YST-2)[%d] = %+v, want %s with the error io.EOF", i, got[i], want.Name)
					}
				} else if got[i] != want {
					t.Errorf("pollAll(YST-1, YST-2)[%d] = %+v, want %+v", i, got[i], want)
				}
				if line := watchLine(got[i]); line != tt.want[i] {
					t.Errorf("watchLine(pollAll(YST-1, YST-2)[%d]) = %q, want %q", i, line, tt.want[i])
				}
			}
		})
	}
}

// The two dishes on Tuesday 25 April: one anemometer, so one wind.
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

// connect serves sim on one end of a pipe and returns a client on the other end that has said
// HELLO. If down is true, the controller hangs up as soon as it is asked for its status, as a
// rebooting one would. When the test ends, both ends close and the simulator stops.
func connect(t *testing.T, sim *Sim, down bool) *Client {
	t.Helper()
	end, simEnd := net.Pipe()
	if down {
		sim.Before = func(req string) {
			if req == "QUERY" {
				simEnd.Close()
			}
		}
	}
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
