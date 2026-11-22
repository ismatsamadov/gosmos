package main

import (
	"iter"
	"slices"
	"testing"
	"time"
)

// benchStream measures assign on Cape Wren's candidates for
// take-over day, read from a stream in AOS order.
func benchStream(b *testing.B) {
	cwr, day := stationDay("CWR", june28)
	slices.SortFunc(day, byAOS)
	next, stop := iter.Pull(slices.Values(day))
	defer stop()
	stream := func(yield func(PassView) bool) {
		for {
			p, ok := next()
			if !ok || !yield(p) {
				return
			}
		}
	}
	for b.Loop() {
		assign(cwr, stream, nil, 2*time.Minute)
	}
}
