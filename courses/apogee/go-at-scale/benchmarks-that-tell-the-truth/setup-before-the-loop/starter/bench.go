package main

import (
	"testing"
	"time"
)

// benchPlanDay measures planDay on Cape Wren's candidates for
// take-over day.
func benchPlanDay(b *testing.B) {
	for b.Loop() {
		cwr, day := stationDay("CWR", june28)
		planDay(day, []Site{cwr}, 2*time.Minute)
	}
}
