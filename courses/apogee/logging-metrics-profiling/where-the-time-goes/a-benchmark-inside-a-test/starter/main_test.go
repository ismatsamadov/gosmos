package main

import (
	"flag"
	"slices"
	"testing"
)

// plan plans in, counting the elevationDeg calls it took.
func plan(in inputs) ([]Planned, int64) {
	before := elevationCalls.Load()
	day := planDay(in.sats, in.els, in.stations, in.bookings, in.from, in.to)
	return day, elevationCalls.Load() - before
}

func TestBenchPlan(t *testing.T) {
	flag.Set("test.benchtime", "100ms")
	want, perPlan := plan(benchInputs())

	setups.Store(0)
	before := elevationCalls.Load()
	r := testing.Benchmark(benchPlan)
	calls := elevationCalls.Load() - before

	if n := setups.Load(); n != 1 {
		t.Errorf("testing.Benchmark(benchPlan) called benchInputs %d times, want 1: before the loop, once", n)
	}
	if calls < 2*perPlan || calls%perPlan != 0 {
		t.Errorf("testing.Benchmark(benchPlan) asked elevationDeg %d times, want a whole number of plans of %d calls each, more than one", calls, perPlan)
	}
	if r.NsPerOp() <= 0 {
		t.Errorf("testing.Benchmark(benchPlan).NsPerOp() = %d, want more than 0", r.NsPerOp())
	}
	if got, _ := plan(benchInputs()); !slices.Equal(got, want) {
		t.Errorf("planDay(benchInputs()) after the benchmark gave %d passes, want the same %d as before", len(got), len(want))
	}
}
