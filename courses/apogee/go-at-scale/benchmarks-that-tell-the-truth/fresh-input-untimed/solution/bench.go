package main

import (
	"slices"
	"testing"
	"time"
)

// benchFresh measures planDay on Cape Wren's candidates for
// take-over day, unsorted, as stationDay predicts them.
func benchFresh(b *testing.B) {
	cwr, day := stationDay("CWR", june28)
	for b.Loop() {
		b.StopTimer()
		fresh := slices.Clone(day)
		b.StartTimer()
		planDay(fresh, []Site{cwr}, 2*time.Minute)
	}
}
