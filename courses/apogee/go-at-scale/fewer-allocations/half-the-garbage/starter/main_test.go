package main

import (
	"flag"
	"runtime"
	"slices"
	"testing"
	"time"
)

// at20 sets -test.benchtime to 20x until the test ends, so that
// testing.Benchmark turns a loop exactly 20 times, and a benchmark run
// from a terminal after the tests keeps the time it asked for.
func at20(t *testing.T) {
	was := flag.Lookup("test.benchtime").Value.String()
	flag.Set("test.benchtime", "20x")
	t.Cleanup(func() { flag.Set("test.benchtime", was) })
}

// garbage returns what one call of plan allocates on s's
// candidates, measured by testing.Benchmark: allocations and bytes.
func garbage(plan func(Site, []PassView, time.Duration) []Contact, s Site, cands []PassView) (allocs, bytes int64) {
	r := testing.Benchmark(func(b *testing.B) {
		for b.Loop() {
			plan(s, cands, 2*time.Minute)
		}
	})
	return r.AllocsPerOp(), r.AllocedBytesPerOp()
}

func TestPlanStationLean(t *testing.T) {
	cwr, cwrDay := stationDay("CWR", june28)
	yst, ystDay := stationDay("YST", june28)
	tests := []struct {
		name  string
		s     Site
		cands []PassView
	}{
		{"Cape Wren's day, one dish", cwr, cwrDay},
		{"Ystmark's day, two dishes", yst, ystDay},
		{"Ystmark's day, every candidate twice", yst, append(slices.Clone(ystDay), ystDay...)},
		{"Cape Wren with no candidates", cwr, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := planStationFat(tt.s, slices.Clone(tt.cands), 2*time.Minute)
			in := slices.Clone(tt.cands)
			got := planStationLean(tt.s, in, 2*time.Minute)
			call := "planStationLean(" + tt.s.Code + ")"
			for i := range min(len(got), len(want)) {
				g, w := got[i], want[i]
				if g.ID != w.ID || g.Dish != w.Dish {
					t.Fatalf("%s[%d] = %s dish %d, want %s dish %d", call, i, g.ID, g.Dish, w.ID, w.Dish)
				}
			}
			if len(got) != len(want) {
				t.Fatalf("%s gave %d contacts, want %d, as planStationFat gives", call, len(got), len(want))
			}
			if !slices.Equal(in, tt.cands) {
				t.Errorf("%s changed its candidates, want them left as they came: planStationFat sorts a copy", call)
			}
		})
	}
}

func TestPlanStationLeanGarbage(t *testing.T) {
	if raceOn {
		t.Skip("allocation counts differ under -race")
	}
	at20(t)
	// One CPU: the pool keeps a part for each CPU, and a
	// benchmark's goroutine may move between them from one turn
	// to the next.
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	for _, code := range []string{"CWR", "YST"} {
		s, day := stationDay(code, june28)
		fatAllocs, fatBytes := garbage(planStationFat, s, day)
		allocs, bytes := garbage(planStationLean, s, day)
		if allocs > fatAllocs/2 {
			t.Errorf("planStationLean(%s's day) made %d allocs/op, want at most %d: half of planStationFat's %d", code, allocs, fatAllocs/2, fatAllocs)
		}
		if bytes > fatBytes/2 {
			t.Errorf("planStationLean(%s's day) allocated %d B/op, want at most %d: half of planStationFat's %d", code, bytes, fatBytes/2, fatBytes)
		}
	}
}
