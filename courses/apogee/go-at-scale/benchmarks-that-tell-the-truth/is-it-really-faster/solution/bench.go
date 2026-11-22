package main

import (
	"fmt"
	"testing"
	"time"
)

// benchWindows measures the planners on windows that start at
// midnight on Friday 18 May.
func benchWindows(b *testing.B) {
	n := newNetwork()
	cfg := Config{MinElevDeg: 10, Turnaround: 2 * time.Minute}
	friday := time.Date(2029, 5, 18, 0, 0, 0, 0, time.UTC)
	planners := []struct {
		name string
		plan func(*Network, Window, Config) []Contact
	}{{"old", planOld}, {"new", planNew}}
	for _, p := range planners {
		b.Run(p.name, func(b *testing.B) {
			for _, hours := range []int{6, 12, 24} {
				to := friday.Add(time.Duration(hours) * time.Hour)
				w := Window{From: friday, To: to}
				b.Run(fmt.Sprintf("hours=%d", hours), func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						p.plan(n, w, cfg)
					}
				})
			}
		})
	}
}
