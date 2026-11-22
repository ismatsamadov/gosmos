package main

import (
	"flag"
	"testing"
	"time"
)

func BenchmarkFresh(b *testing.B) { benchFresh(b) }

// at20 sets -test.benchtime to 20x until the test ends, so that
// testing.Benchmark turns a loop exactly 20 times, and a benchmark run
// from a terminal after the tests keeps the time it asked for.
func at20(t *testing.T) {
	was := flag.Lookup("test.benchtime").Value.String()
	flag.Set("test.benchtime", "20x")
	t.Cleanup(func() { flag.Set("test.benchtime", was) })
}

func TestBenchFresh(t *testing.T) {
	at20(t)
	_, day := stationDay("CWR", june28)

	planned = dayCounts{}
	r := testing.Benchmark(benchFresh)
	got := planned
	if r.N != 20 {
		t.Errorf("testing.Benchmark(benchFresh).N = %d, want 20", r.N)
	}
	if got.calls != 20 {
		t.Fatalf("planDay calls in testing.Benchmark(benchFresh) = %d, want 20: one in each turn of the loop", got.calls)
	}
	if want := 20 * len(day); got.cands != want {
		t.Errorf("candidates handed to planDay = %d, want %d: Cape Wren's %d in each of the 20 calls", got.cands, want, len(day))
	}
	if got.sorted != 0 {
		t.Errorf("planDay calls whose candidates came already sorted = %d, want 0: a fresh copy of the day in every turn", got.sorted)
	}
}

func TestBenchFreshAllocs(t *testing.T) {
	if raceOn {
		t.Skip("allocation counts differ under -race")
	}
	at20(t)
	cwr, day := stationDay("CWR", june28)
	own := int64(testing.AllocsPerRun(5, func() {
		planDay(day, []Site{cwr}, 2*time.Minute)
	}))

	r := testing.Benchmark(benchFresh)
	switch got := r.AllocsPerOp(); {
	case got > own:
		t.Errorf("testing.Benchmark(benchFresh).AllocsPerOp() = %d, want %d, planDay's own: copy the day with the timer stopped", got, own)
	case got < own:
		t.Errorf("testing.Benchmark(benchFresh).AllocsPerOp() = %d, want %d, planDay's own: start the timer again before planDay", got, own)
	}
}
