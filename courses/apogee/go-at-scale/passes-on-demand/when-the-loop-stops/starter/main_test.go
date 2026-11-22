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

func TestPasses(t *testing.T) {
	n := newNetwork()
	sqt, furrowA := n.Sites[2], n.Elements[3]
	july := Window{From: july1, To: july1.AddDate(0, 1, 0)}

	t.Run("Hollis's first three July passes over Salar Quieto", func(t *testing.T) {
		var asked []time.Time
		counting := func(e Elements, s Site, from, to time.Time) []PassView {
			asked = append(asked, from)
			return predict(e, s, from, to)
		}
		got, panicked := upTo(Passes(counting, furrowA, sqt, july), 3)
		if panicked != nil {
			t.Fatalf("breaking after 3 passes panicked: %v", panicked)
		}
		want := predict(furrowA, sqt, july1, july1.AddDate(0, 0, 1))[:3]
		if len(got) != 3 {
			t.Fatalf("Passes(furrow-a, SQT, July) gave the loop %d passes, want 3", len(got))
		}
		for i := range got {
			if got[i].ID != want[i].ID {
				t.Errorf("Passes(furrow-a, SQT, July) pass %d = %s, want %s", i, got[i].ID, want[i].ID)
			}
		}
		for _, from := range asked {
			if !from.Equal(july1) {
				t.Fatalf("after the loop broke at its third pass, predict was asked about %s, want 1 July only", from.Format("2 January"))
			}
		}
	})
	t.Run("every pass of 1 to 7 July", func(t *testing.T) {
		week := Window{From: july1, To: july1.AddDate(0, 0, 7)}
		want := predict(furrowA, sqt, week.From, week.To)
		got, panicked := upTo(Passes(predict, furrowA, sqt, week), len(want)+1)
		if panicked != nil {
			t.Fatalf("Passes(furrow-a, SQT, 1 to 7 July) panicked: %v", panicked)
		}
		if len(got) > len(want) {
			t.Fatalf("Passes(furrow-a, SQT, 1 to 7 July) yielded more than %d passes, want %d", len(want), len(want))
		}
		if len(got) != len(want) {
			t.Fatalf("Passes(furrow-a, SQT, 1 to 7 July) yielded %d passes, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i].ID != want[i].ID {
				t.Fatalf("Passes(furrow-a, SQT, 1 to 7 July) pass %d = %s, want %s", i, got[i].ID, want[i].ID)
			}
		}
	})
	t.Run("an empty window", func(t *testing.T) {
		calls := 0
		counting := func(e Elements, s Site, from, to time.Time) []PassView {
			calls++
			return predict(e, s, from, to)
		}
		empty := Window{From: july1, To: july1}
		got, panicked := upTo(Passes(counting, furrowA, sqt, empty), 1)
		if panicked != nil || len(got) != 0 || calls != 0 {
			t.Errorf("Passes(an empty window) yielded %d passes after %d calls to predict, want 0 and 0", len(got), calls)
		}
	})
}
