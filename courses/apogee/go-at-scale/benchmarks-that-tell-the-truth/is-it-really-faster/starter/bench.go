package main

import (
	"testing"
	"time"
)

// benchWindows measures the planners on windows that start at
// midnight on Friday 18 May.
func benchWindows(b *testing.B) {
	n := newNetwork()
	cfg := Config{MinElevDeg: 10, Turnaround: 2 * time.Minute}
	friday := time.Date(2029, 5, 18, 0, 0, 0, 0, time.UTC)
	day := Window{From: friday, To: friday.Add(24 * time.Hour)}
	for b.Loop() {
		planNew(n, day, cfg)
	}
}
