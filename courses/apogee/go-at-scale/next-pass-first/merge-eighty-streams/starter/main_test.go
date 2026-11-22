package main

import (
	"cmp"
	"iter"
	"math/rand/v2"
	"slices"
	"testing"
	"time"
)

// source yields vals and records in *ended that its function has returned, as a stopped or finished sequence does.
func source(vals []int, ended *bool) iter.Seq[int] {
	return func(yield func(int) bool) {
		defer func() { *ended = true }()
		for _, v := range vals {
			if !yield(v) {
				return
			}
		}
	}
}

// streams returns 80 sorted streams of up to 30 ints each, from a seeded source, and every value in them.
func streams() ([][]int, []int) {
	r := rand.New(rand.NewPCG(1, 2))
	var all []int
	ss := make([][]int, 80)
	for i := range ss {
		for range r.IntN(31) {
			ss[i] = append(ss[i], r.IntN(1000))
		}
		slices.Sort(ss[i])
		all = append(all, ss[i]...)
	}
	slices.Sort(all)
	return ss, all
}

func TestMergeBy(t *testing.T) {
	t.Run("80 sorted streams of ints", func(t *testing.T) {
		ss, want := streams()
		ended := make([]bool, len(ss))
		var seqs []iter.Seq[int]
		for i, s := range ss {
			seqs = append(seqs, source(s, &ended[i]))
		}
		got := slices.Collect(MergeBy(cmp.Compare[int], seqs...))
		if len(got) != len(want) {
			t.Fatalf("MergeBy(80 streams) yielded %d values, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("MergeBy(80 streams) value %d = %d, want %d", i, got[i], want[i])
			}
		}
	})
	t.Run("80 satellites over Cape Wren on Monday 9 April", func(t *testing.T) {
		n := newNetwork()
		cwr := n.Sites[4]
		monday := time.Date(2029, 4, 9, 0, 0, 0, 0, time.UTC)
		w := Window{From: monday, To: monday.AddDate(0, 0, 1)}
		var seqs []iter.Seq[PassView]
		var want []PassView
		for _, e := range n.Elements {
			seqs = append(seqs, Passes(predict, e, cwr, w))
			want = append(want, predict(e, cwr, w.From, w.To)...)
		}
		slices.SortFunc(want, byAOS)
		got := slices.Collect(MergeBy(byAOS, seqs...))
		if len(got) != len(want) {
			t.Fatalf("MergeBy(byAOS, 80 satellites' passes) yielded %d passes, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i].ID != want[i].ID {
				t.Fatalf("MergeBy(byAOS, 80 satellites' passes) pass %d = %s, want %s", i, got[i].ID, want[i].ID)
			}
		}
	})
	t.Run("equal gusts from five stations", func(t *testing.T) {
		type gust struct {
			kmh     int
			station string
		}
		byKmh := func(a, b gust) int { return cmp.Compare(a.kmh, b.kmh) }
		codes := []string{"CWR", "OST", "PVT", "NIM", "VRL"}
		var seqs []iter.Seq[gust]
		var want []string
		for _, code := range codes {
			seqs = append(seqs, slices.Values([]gust{{50, code}, {64, code}, {72, code}}))
		}
		for range 3 {
			want = append(want, codes...)
		}
		var got []string
		for g := range MergeBy(byKmh, seqs...) {
			got = append(got, g.station)
		}
		if len(got) != len(want) {
			t.Fatalf("MergeBy(byKmh, five stations' gusts) yielded %d gusts, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("MergeBy(byKmh, five stations' gusts) gust %d = %s's, want %s's: on a tie, the earlier sequence first", i, got[i], want[i])
			}
		}
	})
	t.Run("a loop that breaks after 10 values", func(t *testing.T) {
		ss, _ := streams()
		ended := make([]bool, len(ss))
		var seqs []iter.Seq[int]
		for i, s := range ss {
			seqs = append(seqs, source(s, &ended[i]))
		}
		taken := 0
		for range MergeBy(cmp.Compare[int], seqs...) {
			if taken++; taken == 10 {
				break
			}
		}
		for i, e := range ended {
			if !e && len(ss[i]) > 0 {
				t.Fatalf("after the loop broke, stream %d was never stopped: call the stop that iter.Pull gave for it", i)
			}
		}
	})
	t.Run("no streams", func(t *testing.T) {
		for v := range MergeBy(cmp.Compare[int]) {
			t.Fatalf("MergeBy(cmp.Compare) yielded %d, want nothing", v)
		}
	})
}
