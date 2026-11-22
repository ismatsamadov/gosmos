package main

import (
	"container/heap"
	"math/rand/v2"
	"slices"
	"testing"
	"time"
)

// candidates returns Monday 2 April's candidates at the ten stations then live: the passes of
// satellites each station can serve that climb to 10° or more, by station.
func candidates() []PassView {
	n := newNetwork()
	day := time.Date(2029, 4, 2, 0, 0, 0, 0, time.UTC)
	var out []PassView
	for _, s := range n.Sites {
		for i, e := range n.Elements {
			if !canServe(s, n.Satellites[i]) {
				continue
			}
			for _, p := range predict(e, s, day, day.AddDate(0, 0, 1)) {
				if p.MaxElevDeg >= 10 {
					out = append(out, p)
				}
			}
		}
	}
	return out
}

// checkPops pops every pass from h and fails t at the first that isn't the next of want.
func checkPops(t *testing.T, h *passHeap, want []PassView) {
	t.Helper()
	for i, w := range want {
		p, _ := heap.Pop(h).(PassView)
		if p.ID != w.ID {
			t.Fatalf("heap.Pop number %d = %q, want %q", i+1, p.ID, w.ID)
		}
	}
	if h.Len() != 0 {
		t.Errorf("Len() after %d pops = %d, want 0", len(want), h.Len())
	}
}

func TestPassHeap(t *testing.T) {
	want := candidates()
	shuffled := slices.Clone(want)
	slices.SortFunc(want, byAOS)
	r := rand.New(rand.NewPCG(1, 2))
	r.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	t.Run("2 April's 286 candidates, pushed in a shuffled order", func(t *testing.T) {
		h := &passHeap{}
		for _, p := range shuffled {
			heap.Push(h, p)
		}
		if h.Len() != len(want) {
			t.Fatalf("Len() after %d calls to heap.Push = %d, want %d", len(want), h.Len(), len(want))
		}
		checkPops(t, h, want)
	})
	t.Run("heap.Init on the same passes, unsorted", func(t *testing.T) {
		h := passHeap(slices.Clone(shuffled))
		heap.Init(&h)
		checkPops(t, &h, want)
	})
	t.Run("a heap of one pass", func(t *testing.T) {
		h := &passHeap{}
		heap.Push(h, want[0])
		if h.Len() != 1 {
			t.Fatalf("Len() after one heap.Push = %d, want 1", h.Len())
		}
		checkPops(t, h, want[:1])
	})
}
