package main

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"testing"
	"time"
)

// schedule is a fixed day of nine passes, out of order, starting an hour from now.
func schedule() []Pass {
	start := time.Now().UTC().Add(time.Hour).Truncate(time.Minute)
	pass := func(sat, station string, aosMin, losMin int) Pass {
		return Pass{SatID: sat, Station: station,
			AOS: start.Add(time.Duration(aosMin) * time.Minute),
			LOS: start.Add(time.Duration(losMin) * time.Minute)}
	}
	return []Pass{
		pass("pipit-1", "YST", 95, 104),
		pass("furrow-a", "SQT", 0, 11),
		pass("furrow-b", "YST", 10, 21),
		pass("cirrine-1", "PEL", 0, 10),
		pass("furrow-a", "YST", 52, 63),
		pass("pipit-3", "CWR", 30, 39),
		pass("pipit-2", "PEL", 61, 70),
		pass("cirrine-1", "YST", 140, 151),
		pass("brine-3", "SQT", 75, 84),
	}
}

// shiftKey marks the context the test gives runDay, so a pass can tell it came from there.
type shiftKey struct{}

// fakeRun stands in for the ground stations. Ystmark's dish stows during furrow-a's pass.
// Port Ellery's and Salar Quieto's first passes share an AOS, so each waits until the other
// has started: run one station after another, they would wait forever.
type fakeRun struct {
	mu       sync.Mutex
	ran      map[string][]string // the satellites each station ran, in order
	problems []string
	pelUp    chan struct{}
	sqtUp    chan struct{}
}

// run is the run function runDay is given.
func (f *fakeRun) run(ctx context.Context, p Pass) error {
	f.mu.Lock()
	f.ran[p.Station] = append(f.ran[p.Station], p.SatID)
	if deadline, ok := ctx.Deadline(); !ok || !deadline.Equal(p.LOS) {
		f.problems = append(f.problems, fmt.Sprintf("%s at %s: ctx.Deadline() = %v, %v; want its LOS %v, true",
			p.SatID, p.Station, deadline, ok, p.LOS))
	}
	if ctx.Value(shiftKey{}) == nil {
		f.problems = append(f.problems, fmt.Sprintf("%s at %s: its ctx doesn't come from the ctx runDay was given",
			p.SatID, p.Station))
	}
	f.mu.Unlock()
	switch {
	case p.Station == "PEL" && p.SatID == "cirrine-1":
		close(f.pelUp)
		<-f.sqtUp
	case p.Station == "SQT" && p.SatID == "furrow-a":
		close(f.sqtUp)
		<-f.pelUp
	case p.Station == "YST" && p.SatID == "furrow-a":
		return fmt.Errorf("YST: wind 81 km/h: %w", ErrStowed)
	}
	return nil
}

func TestRunDay(t *testing.T) {
	f := &fakeRun{ran: map[string][]string{}, pelUp: make(chan struct{}), sqtUp: make(chan struct{})}
	result := make(chan Report, 1)
	shift := context.WithValue(context.Background(), shiftKey{}, "friday")
	go func() {
		result <- runDay(shift, schedule(), f.run)
	}()
	var got Report
	select {
	case got = <-result:
	case <-time.After(time.Second):
		t.Fatal("runDay(...) is still running after 1s: PEL and SQT have passes at the same AOS, so run the stations at the same time, holding no lock while a station runs")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if want := (Report{Done: 6, Failed: 1, Skipped: 2}); got != want {
		t.Errorf("runDay(...) = %+v, want %+v", got, want)
	}
	wantRan := map[string][]string{
		"YST": {"furrow-b", "furrow-a"}, // stowed during furrow-a: the last two never run
		"PEL": {"cirrine-1", "pipit-2"},
		"SQT": {"furrow-a", "brine-3"},
		"CWR": {"pipit-3"},
	}
	for _, station := range []string{"YST", "PEL", "SQT", "CWR"} {
		if got, want := f.ran[station], wantRan[station]; !slices.Equal(got, want) {
			t.Errorf("runDay(...) ran %v at %s, want %v, in that order", got, station, want)
		}
	}
	if len(f.problems) > 0 {
		t.Errorf("%d problems with the passes' contexts; the first: %s", len(f.problems), f.problems[0])
	}
}
