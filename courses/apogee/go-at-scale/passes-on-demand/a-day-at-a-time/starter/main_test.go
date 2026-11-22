package main

import (
	"iter"
	"testing"
	"time"
)

// july1 is the first day of the third quarter, Hollis's summer.
var july1 = time.Date(2029, 7, 1, 0, 0, 0, 0, time.UTC)

// upTo ranges over seq until it ends or has given limit passes, and returns them with
// whatever the loop panicked with.
func upTo(seq iter.Seq[PassView], limit int) (got []PassView, panicked any) {
	defer func() { panicked = recover() }()
	for p := range seq {
		got = append(got, p)
		if len(got) == limit {
			break
		}
	}
	return got, nil
}

func TestCatalog(t *testing.T) {
	n := newNetwork()
	week := Window{From: july1, To: july1.AddDate(0, 0, 7)}

	t.Run("the first week of the third quarter", func(t *testing.T) {
		const want = 51174
		got, panicked := upTo(Catalog(n, week, predict), want+1)
		if panicked != nil {
			t.Fatalf("Catalog(1 to 7 July) panicked: %v", panicked)
		}
		if len(got) > want {
			t.Fatalf("Catalog(1 to 7 July) yielded more than %d passes, want %d", want, want)
		}
		if len(got) != want {
			t.Fatalf("Catalog(1 to 7 July) yielded %d passes, want %d", len(got), want)
		}
		for i, id := range []string{"brine-27/VSK/00697", "furrow-s/SRH/00049", "brine-16/CRP/00697"} {
			if got[i].ID != id {
				t.Fatalf("Catalog(1 to 7 July) pass %d = %s, want %s", i, got[i].ID, id)
			}
		}
		sats, sites := map[string]Satellite{}, map[string]Site{}
		for _, sat := range n.Satellites {
			sats[sat.ID] = sat
		}
		for _, s := range n.Sites {
			sites[s.Code] = s
		}
		for i, p := range got {
			if !canServe(sites[p.Station], sats[p.SatID]) {
				t.Fatalf("Catalog(1 to 7 July) yielded %s, but %s can't serve %s", p.ID, p.Station, p.SatID)
			}
			if i > 0 && p.AOS.Before(got[i-1].AOS) {
				t.Fatalf("Catalog(1 to 7 July) yielded %s before %s, want AOS order", got[i-1].ID, p.ID)
			}
		}
	})
	t.Run("Hollis stops after ten passes", func(t *testing.T) {
		july2 := july1.AddDate(0, 0, 1)
		var askedFrom, askedTo time.Time // the first call about more than 1 July
		counting := func(e Elements, s Site, from, to time.Time) []PassView {
			if from.Before(july1) || to.After(july2) {
				if askedFrom.IsZero() {
					askedFrom, askedTo = from, to
				}
				return nil // the test fails on this call: predicting it would only take time
			}
			return predict(e, s, from, to)
		}
		got, panicked := upTo(Catalog(n, week, counting), 10)
		if panicked != nil {
			t.Fatalf("breaking after 10 passes panicked: %v", panicked)
		}
		if !askedFrom.IsZero() {
			t.Fatalf("Catalog(1 to 7 July) asked predict about %s to %s for a loop that stops at its tenth pass, want 1 July only", askedFrom.Format("2 January 15:04"), askedTo.Format("2 January 15:04"))
		}
		if len(got) != 10 {
			t.Fatalf("Catalog(1 to 7 July) gave the loop %d passes, want 10", len(got))
		}
	})
}
