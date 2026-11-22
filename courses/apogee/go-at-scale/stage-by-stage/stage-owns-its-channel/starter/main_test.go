package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

// portEllery is three of Cirrine-1's frames as Port Ellery heard them on Monday 5 March.
var portEllery = []string{
	"cirrine-1 3202 2029-03-05T09:25:05Z batt_v=8.05 temp_c=10.8 mode=NOMINAL",
	"cirrine-1 3203 2029-03-05T09:25:35Z batt_v=8.05 temp_c=10.9 mode=NOMINAL",
	"cirrine-1 3204 2029-03-05T09:26:05Z batt_v=8.04 temp_c=10.9 mode=NOMINAL",
}

// frames returns n of sat's telemetry lines, one every 30 seconds from at, numbered from seq.
func frames(sat string, seq int, at time.Time, n int) []string {
	lines := make([]string, n)
	for i := range lines {
		t := at.Add(time.Duration(i) * 30 * time.Second)
		lines[i] = fmt.Sprintf("%s %d %s batt_v=8.01 temp_c=10.5 mode=NOMINAL", sat, seq+i, t.Format(time.RFC3339))
	}
	return lines
}

// guard is how long a test waits for what a right stage does at once: one second in all.
type guard struct{ end time.Time }

// newGuard starts a guard's second.
func newGuard() guard { return guard{end: time.Now().Add(time.Second)} }

// wait is ready when the guard's second is over, and never sooner than 50 ms from now.
func (g guard) wait() <-chan time.Time {
	return time.After(max(time.Until(g.end), 50*time.Millisecond))
}

func TestRead(t *testing.T) {
	g := newGuard()
	night := time.Date(2029, 3, 5, 5, 20, 55, 0, time.UTC)
	ystmark := append(frames("pipit-2", 5510, night, 20), frames("cirrine-1", 2991, night.Add(2*time.Hour+18*time.Minute+40*time.Second), 20)...)
	tests := []struct {
		name, station string
		lines         []string
	}{
		{"Port Ellery's three lines", "PEL", portEllery},
		{"Tirrin Rock's feed, empty until April", "TIR", nil},
		{"Ystmark's 40 lines, more than the channel's room", "YST", ystmark},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := read(context.Background(), tt.station, strings.NewReader(strings.Join(tt.lines, "\n")))
			var got []Heard
			for open := true; open; {
				select {
				case h, ok := <-ch:
					if ok {
						got = append(got, h)
					}
					open = ok
				case <-g.wait():
					t.Fatalf("read did not close its channel: %d of %d lines from %s received, then nothing", len(got), len(tt.lines), tt.station)
				}
			}
			if len(got) != len(tt.lines) {
				t.Fatalf("read(ctx, %q, feed) sent %d lines, want %d", tt.station, len(got), len(tt.lines))
			}
			for i, h := range got {
				if want := (Heard{Station: tt.station, Line: tt.lines[i]}); h != want {
					t.Errorf("read(ctx, %q, feed) value %d = {%q, %q}, want {%q, %q}", tt.station, i, h.Station, h.Line, want.Station, want.Line)
				}
			}
		})
	}
}
