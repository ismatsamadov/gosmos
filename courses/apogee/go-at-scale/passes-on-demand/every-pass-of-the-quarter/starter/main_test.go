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

func TestBookable(t *testing.T) {
	n := newNetwork()
	twoWeeks := Window{From: july1, To: july1.AddDate(0, 0, 14)}
	sats, sites := map[string]Satellite{}, map[string]Site{}
	for _, sat := range n.Satellites {
		sats[sat.ID] = sat
	}
	for _, s := range n.Sites {
		sites[s.Code] = s
	}

	t.Run("Furrowlight's first two weeks of July", func(t *testing.T) {
		asked := map[string]bool{}
		counting := func(e Elements, s Site, from, to time.Time) []PassView {
			asked[e.SatID] = true
			if sats[e.SatID].Operator != "Furrowlight" {
				return nil // the test fails on this call: predicting it would only take time
			}
			return predict(e, s, from, to)
		}
		const want = 37917
		got, panicked := upTo(Bookable(n, twoWeeks, counting, "Furrowlight"), want+1)
		if panicked != nil {
			t.Fatalf("Bookable(Furrowlight, 1 to 14 July) panicked: %v", panicked)
		}
		for _, sat := range n.Satellites {
			if asked[sat.ID] && sat.Operator != "Furrowlight" {
				t.Fatalf("Bookable(Furrowlight, 1 to 14 July) asked predict about %s, which is %s's", sat.ID, sat.Operator)
			}
		}
		if len(got) > want {
			t.Fatalf("Bookable(Furrowlight, 1 to 14 July) yielded more than %d passes, want %d", want, want)
		}
		if len(got) != want {
			t.Fatalf("Bookable(Furrowlight, 1 to 14 July) yielded %d passes, want %d", len(got), want)
		}
		if got[0].ID != "furrow-s/SRH/00049" {
			t.Fatalf("Bookable(Furrowlight, 1 to 14 July) pass 0 = %s, want furrow-s/SRH/00049", got[0].ID)
		}
		for i, p := range got {
			if sat := sats[p.SatID]; sat.Operator != "Furrowlight" || !canServe(sites[p.Station], sat) {
				t.Fatalf("Bookable(Furrowlight, 1 to 14 July) yielded %s, not a pass of Furrowlight's at a station that can serve it", p.ID)
			}
			if i > 0 && p.AOS.Before(got[i-1].AOS) {
				t.Fatalf("Bookable(Furrowlight, 1 to 14 July) yielded %s before %s, want AOS order", got[i-1].ID, p.ID)
			}
		}
	})
	t.Run("a request with no operator", func(t *testing.T) {
		calls := 0
		counting := func(e Elements, s Site, from, to time.Time) []PassView {
			calls++
			return nil // any call fails the test: predicting it would only take time
		}
		got, panicked := upTo(Bookable(n, twoWeeks, counting, ""), 1)
		if panicked != nil || len(got) != 0 || calls != 0 {
			t.Errorf("Bookable(no operator, 1 to 14 July) yielded %d passes after %d calls to predict, want 0 and 0", len(got), calls)
		}
	})
	t.Run("Hollis's program stops after 20 passes", func(t *testing.T) {
		type call struct{ from, to time.Time }
		var calls []call
		counting := func(e Elements, s Site, from, to time.Time) []PassView {
			calls = append(calls, call{from, to})
			return predict(e, s, from, to)
		}
		got, panicked := upTo(Bookable(n, twoWeeks, counting, "Furrowlight"), 20)
		if panicked != nil {
			t.Fatalf("breaking after 20 passes panicked: %v", panicked)
		}
		if len(got) != 20 {
			t.Fatalf("Bookable(Furrowlight, 1 to 14 July) gave the loop %d passes, want 20", len(got))
		}
		july2 := july1.AddDate(0, 0, 1)
		for _, c := range calls {
			if c.from.Before(july1) || c.to.After(july2) {
				t.Fatalf("after the loop broke at its twentieth pass, predict had been asked about %s to %s, want 1 July only", c.from.Format("2 January 15:04"), c.to.Format("2 January 15:04"))
			}
		}
	})
}
