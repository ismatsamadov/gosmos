package main

import (
	"flag"
	"testing"
	"time"
)

func BenchmarkPlanDay(b *testing.B) { benchPlanDay(b) }

// at20 sets -test.benchtime to 20x until the test ends, so that
// testing.Benchmark turns a loop exactly 20 times, and a benchmark run
// from a terminal after the tests keeps the time it asked for.
func at20(t *testing.T) {
	was := flag.Lookup("test.benchtime").Value.String()
	flag.Set("test.benchtime", "20x")
	t.Cleanup(func() { flag.Set("test.benchtime", was) })
}

func TestBenchPlanDay(t *testing.T) {
	at20(t)
	cwr, day := stationDay("CWR", june28)
	contacts := len(planDay(day, []Site{cwr}, 2*time.Minute))

	planned = dayCounts{}
	r := testing.Benchmark(benchPlanDay)
	got := planned
	if r.N != 20 {
		t.Errorf("testing.Benchmark(benchPlanDay).N = %d, want 20", r.N)
	}
	if got.calls != 20 {
		t.Fatalf("planDay calls in testing.Benchmark(benchPlanDay) = %d, want 20: one in each turn of the loop", got.calls)
	}
	if want := 20 * len(day); got.cands != want {
		t.Errorf("candidates handed to planDay = %d, want %d: Cape Wren's %d in each of the 20 calls", got.cands, want, len(day))
	}
	if want := 20 * contacts; got.contacts != want {
		t.Errorf("contacts planDay returned = %d, want %d: Cape Wren's %d in each of the 20 calls", got.contacts, want, contacts)
	}
}
