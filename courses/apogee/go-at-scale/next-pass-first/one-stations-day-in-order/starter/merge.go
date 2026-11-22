package main

import "iter"

// head is one sequence's next value in a merge, with the index of that sequence.
type head[T any] struct {
	v T
	i int
}

// byHead orders heads by cmp, and a tie by index, so the earlier sequence's value comes first.
func byHead[T any](cmp func(a, b T) int) func(a, b head[T]) int {
	return func(a, b head[T]) int {
		if c := cmp(a.v, b.v); c != 0 {
			return c
		}
		return a.i - b.i
	}
}

// MergeBy merges sequences sorted by cmp into one sorted sequence, holding one value of each at a time.
func MergeBy[T any](cmp func(a, b T) int, seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		q := NewQueue(byHead(cmp))
		nexts := make([]func() (T, bool), len(seqs))
		for i, seq := range seqs {
			next, stop := iter.Pull(seq)
			defer stop()
			nexts[i] = next
			if v, ok := next(); ok {
				q.Push(head[T]{v, i})
			}
		}
		for {
			h, ok := q.Pop()
			if !ok || !yield(h.v) {
				return
			}
			if v, ok := nexts[h.i](); ok {
				q.Push(head[T]{v, h.i})
			}
		}
	}
}

// stationPasses yields s's candidates in w by AOS, then ID: the passes of every satellite s can serve that climb to minElevDeg or more, predicted a day at a time.
func stationPasses(n *Network, s Site, w Window, predict Predictor, minElevDeg float64) iter.Seq[PassView] {
	return func(yield func(PassView) bool) {}
}
