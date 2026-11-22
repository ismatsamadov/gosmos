package main

import (
	"slices"
	"testing"
	"time"
)

// benchStream measures assign on Cape Wren's candidates for
// take-over day, read from a stream in AOS order.
func benchStream(b *testing.B) {
	cwr, day := stationDay("CWR", june28)
	slices.SortFunc(day, byAOS)
	contacts := 0
	for b.Loop() {
		c, _ := assign(cwr, slices.Values(day), nil, 2*time.Minute)
		contacts += len(c)
	}
	b.ReportMetric(float64(contacts)/float64(b.N), "contacts/op")
}
