package main

import (
	"flag"
	"slices"
	"testing"
	"time"
)

func BenchmarkStream(b *testing.B) { benchStream(b) }

// at20 sets -test.benchtime to 20x until the test ends, so that
// testing.Benchmark turns a loop exactly 20 times, and a benchmark run
// from a terminal after the tests keeps the time it asked for.
func at20(t *testing.T) {
	was := flag.Lookup("test.benchtime").Value.String()
	flag.Set("test.benchtime", "20x")
	t.Cleanup(func() { flag.Set("test.benchtime", was) })
}

func TestBenchStream(t *testing.T) {
	at20(t)
	cwr, day := stationDay("CWR", june28)
	slices.SortFunc(day, byAOS)
	contacts, _ := assign(cwr, slices.Values(day), nil, 2*time.Minute)
	want := float64(len(contacts))

	r := testing.Benchmark(benchStream)
	if r.N != 20 {
		t.Errorf("testing.Benchmark(benchStream).N = %d, want 20", r.N)
	}
	got, ok := r.Extra["contacts/op"]
	if !ok {
		t.Fatalf("testing.Benchmark(benchStream) reported no contacts/op, want %v: b.ReportMetric after the loop", want)
	}
	if got != want {
		t.Errorf("testing.Benchmark(benchStream) reported %v contacts/op, want %v: Cape Wren's contacts, in every turn", got, want)
	}
}

func TestBenchStreamAllocs(t *testing.T) {
	if raceOn {
		t.Skip("allocation counts differ under -race")
	}
	at20(t)
	cwr, day := stationDay("CWR", june28)
	slices.SortFunc(day, byAOS)
	own := int64(testing.AllocsPerRun(5, func() {
		assign(cwr, slices.Values(day), nil, 2*time.Minute)
	}))

	r := testing.Benchmark(benchStream)
	if got := r.AllocsPerOp(); got < own || got > own+2 {
		t.Errorf("testing.Benchmark(benchStream).AllocsPerOp() = %d, want %d to %d: assign's own for the whole day, in every turn", got, own, own+2)
	}
}
