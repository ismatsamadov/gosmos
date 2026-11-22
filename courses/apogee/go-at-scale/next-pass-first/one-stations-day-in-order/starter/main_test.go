package main

import (
	"slices"
	"testing"
	"time"
)

// candidates returns s's candidates in w, predicted all at once and sorted: what stationPasses must yield.
func candidates(n *Network, s Site, w Window) []PassView {
	var out []PassView
	for i, e := range n.Elements {
		if !canServe(s, n.Satellites[i]) {
			continue
		}
		for _, p := range predict(e, s, w.From, w.To) {
			if p.MaxElevDeg >= 10 {
				out = append(out, p)
			}
		}
	}
	slices.SortFunc(out, byAOS)
	return out
}

func TestStationPasses(t *testing.T) {
	n := newNetwork()
	tir, vrl := n.Sites[5], n.Sites[19]
	tuesday := time.Date(2029, 4, 17, 0, 0, 0, 0, time.UTC)
	day := Window{From: tuesday, To: tuesday.AddDate(0, 0, 1)}

	t.Run("Tirrin Rock on Tuesday 17 April", func(t *testing.T) {
		want := candidates(n, tir, day)
		got := slices.Collect(stationPasses(n, tir, day, predict, 10))
		if len(got) != len(want) {
			t.Fatalf("stationPasses(TIR, 17 April) yielded %d candidates, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i].ID != want[i].ID {
				t.Fatalf("stationPasses(TIR, 17 April) candidate %d = %s, want %s", i, got[i].ID, want[i].ID)
			}
		}
		if contacts, dropped := assign(tir, slices.Values(got), nil, 2*time.Minute); len(contacts) != 120 || dropped != 28 {
			t.Errorf("assign(TIR, those candidates) = %d contacts, %d dropped, want 120 and 28", len(contacts), dropped)
		}
	})
	t.Run("Vorrel Island, not live until 4 June", func(t *testing.T) {
		for p := range stationPasses(n, vrl, day, predict, 10) {
			t.Fatalf("stationPasses(VRL, 17 April) yielded %s, want nothing", p.ID)
		}
	})
	t.Run("a loop that stops at the first candidate of the week", func(t *testing.T) {
		week := Window{From: tuesday, To: tuesday.AddDate(0, 0, 7)}
		type call struct {
			day  time.Time
			high bool // the day gave the satellite a candidate
		}
		calls := map[string][]call{}
		counting := func(e Elements, s Site, from, to time.Time) []PassView {
			ps := predict(e, s, from, to)
			high := slices.ContainsFunc(ps, func(p PassView) bool { return p.MaxElevDeg >= 10 })
			calls[e.SatID] = append(calls[e.SatID], call{from, high})
			return ps
		}
		for range stationPasses(n, tir, week, counting, 10) {
			break
		}
		for _, e := range n.Elements {
			cs := calls[e.SatID]
			for k := 1; k < len(cs); k++ {
				if cs[k-1].high {
					t.Fatalf("after one candidate, predict was asked about %s on %s, though %s had already given it a candidate", e.SatID, cs[k].day.Format("2 January"), cs[k-1].day.Format("2 January"))
				}
			}
		}
	})
}
