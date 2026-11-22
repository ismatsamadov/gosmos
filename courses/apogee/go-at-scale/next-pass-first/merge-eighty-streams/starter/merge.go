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
		// TODO: merge through a Queue of heads, not one sequence after another
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}
