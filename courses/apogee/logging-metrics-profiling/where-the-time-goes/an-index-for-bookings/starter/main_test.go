package main

import (
	"flag"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// fortnight is every pass passplan predicts over the five stations for the fourteen days from
// Tuesday 16 January, 00:00, in AOS order.
var fortnight = sync.OnceValue(func() []PassView {
	els := elements()
	var passes []PassView
	for _, sat := range fleet() {
		elems := els[sat.ID]
		for _, st := range allStations() {
			el := func(at time.Time) float64 { return elevationDeg(sat, elems, st, at) }
			passes = append(passes, passViews(sat, elems, st, findPasses(el, epoch, epoch.AddDate(0, 0, 14), st.MaskDeg))...)
		}
	}
	slices.SortFunc(passes, func(a, b PassView) int {
		if c := a.AOS.Compare(b.AOS); c != 0 {
			return c
		}
		return strings.Compare(a.ID, b.ID)
	})
	return passes
})

// fortnightBookings returns Tuesday's bookings, then one for most of the later passes, 2,000 in
// all: passplan read them all for every pass. Speed depends only on their satellites, stations and
// times, so the later ones carry no IDs.
func fortnightBookings(passes []PassView) []Booking {
	bookings := tuesdayBookings(passes)
	for i, p := range passes[300:] {
		if len(bookings) < 2000 && i%6 != 5 {
			bookings = append(bookings, Booking{SatID: p.SatID, Station: p.Station, AOS: p.AOS, LOS: p.LOS})
		}
	}
	return bookings
}

// find returns the pass of the fortnight with id.
func find(t *testing.T, id string) PassView {
	t.Helper()
	for _, p := range fortnight() {
		if p.ID == id {
			return p
		}
	}
	t.Fatalf("no pass %s in the fortnight", id)
	return PassView{}
}

func TestAttach(t *testing.T) {
	passes := fortnight()[:300]
	bookings := fortnightBookings(fortnight())
	got, want := attach(passes, bookings), scanAllBookings(passes, bookings)
	if len(got) != len(want) {
		t.Fatalf("attach(the fortnight's first 300 passes, 2,000 bookings) gave %d passes, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("attach(...)[%d] = %s booked %q, want %s booked %q", i, got[i].ID, got[i].Booking, want[i].ID, want[i].Booking)
		}
	}
}

func TestAttachOnePass(t *testing.T) {
	brine := find(t, "brine-3/CWR/15442")
	lan := find(t, "brine-3/LAN/15447")
	booking := func(id string, p PassView) Booking {
		return Booking{ID: id, SatID: p.SatID, Station: p.Station, AOS: p.AOS, LOS: p.LOS}
	}
	atoll := booking("bk-5799", lan)
	atoll.AOS = atoll.AOS.In(time.FixedZone("+0530", 5*3600+30*60))
	tests := []struct {
		name     string
		pass     PassView
		bookings []Booking
		want     string
	}{
		{"a pass nobody booked", find(t, "furrow-a/LAN/22054"), tuesdayBookings(fortnight()), ""},
		{"a pass booked twice keeps the first booking", brine, []Booking{booking("bk-5803", brine), booking("bk-5805", brine)}, "bk-5803"},
		{"a booking whose AOS is in Lantern Atoll's time, +05:30", lan, []Booking{atoll}, "bk-5799"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := attach([]PassView{tt.pass}, tt.bookings)
			if len(got) != 1 {
				t.Fatalf("attach([%s], bookings) gave %d passes, want 1", tt.pass.ID, len(got))
			}
			if got[0].PassView != tt.pass || got[0].Booking != tt.want {
				t.Errorf("attach([%s], bookings)[0] = %s booked %q, want %s booked %q", tt.pass.ID, got[0].ID, got[0].Booking, tt.pass.ID, tt.want)
			}
		})
	}
}

func TestAttachSpeed(t *testing.T) {
	flag.Set("test.benchtime", "100ms")
	passes := fortnight()
	bookings := fortnightBookings(passes)
	scan := testing.Benchmark(func(b *testing.B) {
		for b.Loop() {
			scanAllBookings(passes, bookings)
		}
	})
	index := testing.Benchmark(func(b *testing.B) {
		for b.Loop() {
			attach(passes, bookings)
		}
	})
	if index.NsPerOp()*5 > scan.NsPerOp() {
		t.Errorf("attach(the fortnight's %d passes, %d bookings) took %d ns a run, want at most %d, a fifth of scanAllBookings's %d", len(passes), len(bookings), index.NsPerOp(), scan.NsPerOp()/5, scan.NsPerOp())
	}
}
