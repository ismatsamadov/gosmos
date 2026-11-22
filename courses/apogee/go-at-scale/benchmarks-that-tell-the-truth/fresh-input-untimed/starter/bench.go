package main

import (
	"testing"
	"time"
)

// benchFresh measures planDay on Cape Wren's candidates for
// take-over day, unsorted, as stationDay predicts them.
func benchFresh(b *testing.B) {
	cwr, day := stationDay("CWR", june28)
	for b.Loop() {
		planDay(day, []Site{cwr}, 2*time.Minute)
	}
}
