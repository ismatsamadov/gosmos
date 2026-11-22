package main

import (
	"flag"
	"io"
	"os"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
)

// oldPlanDay is planDay as passplan had it on Tuesday, kept to compare with.
func oldPlanDay(sats []Satellite, els map[string]Elements, stations []Station, bookings []Booking, from, to time.Time) []Planned {
	var groups [][]Planned
	for _, sat := range sats {
		elems := els[sat.ID]
		for _, st := range stations {
			el := elevation(sat, elems, st)
			passes := passViews(sat, elems, st, scanEverySecond(el, from, to, st.MaskDeg))
			groups = append(groups, scanAllBookings(passes, bookings))
		}
	}
	return sortAfterEach(groups)
}

// day is what the tests plan: Brine-3, Cirrine-1 and Furrow-A over Cape Wren and Lantern Atoll for
// twelve hours from 05:10 on Tuesday 16 January, with that day's bookings.
type day struct {
	sats     []Satellite
	stations []Station
	bookings []Booking
	from, to time.Time
}

// tuesday builds the day once, booking its passes as Tuesday's bookings did.
var tuesday = sync.OnceValue(func() day {
	d := day{
		sats:     []Satellite{satellite("brine-3"), satellite("cirrine-1"), satellite("furrow-a")},
		stations: []Station{station("CWR"), station("LAN")},
		from:     time.Date(2029, 1, 16, 5, 10, 0, 0, time.UTC),
	}
	d.to = d.from.Add(12 * time.Hour)
	var passes []PassView
	for _, p := range oldPlanDay(d.sats, elements(), d.stations, nil, d.from, d.to) {
		passes = append(passes, p.PassView)
	}
	d.bookings = tuesdayBookings(passes)
	return d
})

// plan plans d with plan, counting the elevationDeg calls it took.
func (d day) plan(plan func([]Satellite, map[string]Elements, []Station, []Booking, time.Time, time.Time) []Planned) ([]Planned, int64) {
	before := elevationCalls.Load()
	p := plan(d.sats, elements(), d.stations, d.bookings, d.from, d.to)
	return p, elevationCalls.Load() - before
}

func TestPlanDay(t *testing.T) {
	want, oldCalls := tuesday().plan(oldPlanDay)
	got, calls := tuesday().plan(planDay)
	if len(got) != len(want) {
		t.Fatalf("planDay(Tuesday, 05:10 to 17:10) gave %d passes, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("planDay(Tuesday, 05:10 to 17:10)[%d] = %s booked %q, want %s booked %q", i, got[i].ID, got[i].Booking, want[i].ID, want[i].Booking)
		}
	}
	if calls*20 > oldCalls {
		t.Errorf("planDay(Tuesday, 05:10 to 17:10) asked elevationDeg %d times, want at most %d, a twentieth of the old planner's %d", calls, oldCalls/20, oldCalls)
	}
}

func TestPlanDaySpeed(t *testing.T) {
	flag.Set("test.benchtime", "100ms")
	d := tuesday()
	old := testing.Benchmark(func(b *testing.B) {
		for b.Loop() {
			d.plan(oldPlanDay)
		}
	})
	now := testing.Benchmark(func(b *testing.B) {
		for b.Loop() {
			d.plan(planDay)
		}
	})
	if now.NsPerOp()*5 > old.NsPerOp() {
		t.Errorf("planDay(Tuesday, 05:10 to 17:10) took %d ns a run, want at most %d, a fifth of the old planner's %d", now.NsPerOp(), old.NsPerOp()/5, old.NsPerOp())
	}
}

func TestPlanDayProfile(t *testing.T) {
	flag.Set("test.benchtime", "100ms")
	os.Remove("cpu-after.prof")
	d := tuesday()
	work := func() { // plans Tuesday again and again for about 100 ms, as a benchmark does
		testing.Benchmark(func(b *testing.B) {
			for b.Loop() {
				d.plan(planDay)
			}
		})
	}
	if err := profileTo("cpu-after.prof", work); err != nil {
		t.Fatalf(`profileTo("cpu-after.prof", plan Tuesday for 100 ms) = %v, want <nil>`, err)
	}
	b, err := os.ReadFile("cpu-after.prof")
	switch {
	case err != nil:
		t.Errorf("reading cpu-after.prof gives %v, want the profile", err)
	case len(b) < 2 || b[0] != 0x1f || b[1] != 0x8b:
		t.Errorf("cpu-after.prof starts with % x, want 1f 8b, the first bytes of a profile", b[:min(len(b), 2)])
	}
	if err := pprof.StartCPUProfile(io.Discard); err != nil {
		t.Errorf("pprof.StartCPUProfile(io.Discard) after profileTo = %v, want <nil>", err)
		return
	}
	pprof.StopCPUProfile()
}
