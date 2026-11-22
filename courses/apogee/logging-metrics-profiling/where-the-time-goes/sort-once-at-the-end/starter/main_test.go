package main

import (
	"flag"
	"slices"
	"strings"
	"testing"
	"time"
)

// firstGroups returns the first 600 passes passplan predicts from Tuesday 16 January, 00:00, in the
// groups it builds them in: one for each satellite and station, in AOS order.
func firstGroups() [][]Planned {
	els := elements()
	var groups [][]Planned
	var all []time.Time
	for _, sat := range fleet() {
		elems := els[sat.ID]
		for _, st := range allStations() {
			el := func(at time.Time) float64 { return elevationDeg(sat, elems, st, at) }
			var g []Planned
			for _, p := range passViews(sat, elems, st, findPasses(el, epoch, epoch.AddDate(0, 0, 4), st.MaskDeg)) {
				g = append(g, Planned{PassView: p})
				all = append(all, p.AOS)
			}
			groups = append(groups, g)
		}
	}
	slices.SortFunc(all, time.Time.Compare)
	for i, g := range groups {
		groups[i] = slices.DeleteFunc(g, func(p Planned) bool { return p.AOS.After(all[599]) })
	}
	return groups
}

// copyOf returns a copy of groups that shares nothing with them.
func copyOf(groups [][]Planned) [][]Planned {
	c := make([][]Planned, len(groups))
	for i, g := range groups {
		c[i] = slices.Clone(g)
	}
	return c
}

func TestBuildDay(t *testing.T) {
	groups := firstGroups()
	before := copyOf(groups)
	want := sortAfterEach(copyOf(groups))
	got := buildDay(groups)
	if len(got) != len(want) {
		t.Fatalf("buildDay(600 passes in %d groups) gave %d passes, want %d", len(groups), len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("buildDay(600 passes)[%d] = %s at %s, want %s at %s", i, got[i].ID, got[i].AOS.Format(time.RFC3339), want[i].ID, want[i].AOS.Format(time.RFC3339))
		}
	}
	for i := range groups {
		if !slices.Equal(groups[i], before[i]) {
			t.Errorf("buildDay(groups) changed group %d, want the groups left as they were", i)
			break
		}
	}
}

// TestBuildDayRoomToSpare gives buildDay a first group with room to grow: two of Brine-3's passes
// over Cape Wren, then Cirrine-1's earlier one. The first group must stay as it was.
func TestBuildDayRoomToSpare(t *testing.T) {
	var brine, cirrine []Planned
	for _, g := range firstGroups() {
		switch {
		case len(g) == 0:
		case strings.HasPrefix(g[0].ID, "brine-3/CWR/"):
			brine = slices.Grow(slices.Clone(g[1:3]), 8)
		case strings.HasPrefix(g[0].ID, "cirrine-1/CWR/"):
			cirrine = g[:1]
		}
	}
	before := slices.Clone(brine)
	buildDay([][]Planned{brine, cirrine})
	if !slices.Equal(brine, before) {
		t.Errorf("buildDay([[%s %s] [%s]]) changed the first group to [%s %s], want it left as it was", before[0].ID, before[1].ID, cirrine[0].ID, brine[0].ID, brine[1].ID)
	}
}

// TestBuildDayTies plans the two passes that rise in the same second on Saturday 20 January:
// Furrow-B over Port Ellery and Brine-3 over Lantern Atoll, at 22:47:35.
func TestBuildDayTies(t *testing.T) {
	at := func(s string) time.Time {
		v, _ := time.Parse(time.RFC3339, s)
		return v
	}
	furrow := Planned{PassView: PassView{ID: "furrow-b/PEL/22119", Pass: Pass{SatID: "furrow-b", Station: "PEL", AOS: at("2029-01-20T22:47:35Z"), LOS: at("2029-01-20T22:55:12Z"), MaxElevDeg: 18.0}}}
	brine := Planned{PassView: PassView{ID: "brine-3/LAN/15514", Pass: Pass{SatID: "brine-3", Station: "LAN", AOS: at("2029-01-20T22:47:35Z"), LOS: at("2029-01-20T22:55:42Z"), MaxElevDeg: 88.3}}}
	got := buildDay([][]Planned{{furrow}, {brine}})
	if len(got) != 2 || got[0].ID != brine.ID || got[1].ID != furrow.ID {
		var ids []string
		for _, p := range got {
			ids = append(ids, p.ID)
		}
		t.Errorf("buildDay([[%s] [%s]]) = %v, want [%s %s]: equal AOS, so by ID", furrow.ID, brine.ID, ids, brine.ID, furrow.ID)
	}
}

func TestBuildDaySpeed(t *testing.T) {
	flag.Set("test.benchtime", "100ms")
	groups := firstGroups()
	each := testing.Benchmark(func(b *testing.B) {
		for b.Loop() {
			sortAfterEach(groups)
		}
	})
	once := testing.Benchmark(func(b *testing.B) {
		for b.Loop() {
			buildDay(groups)
		}
	})
	if once.NsPerOp()*5 > each.NsPerOp() {
		t.Errorf("buildDay(600 passes) took %d ns a run, want at most %d, a fifth of sortAfterEach's %d", once.NsPerOp(), each.NsPerOp()/5, each.NsPerOp())
	}
}
