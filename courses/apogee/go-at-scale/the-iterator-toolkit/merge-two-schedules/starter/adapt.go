package main

import "iter"

// Filter yields the values of seq that keep keeps, and stops seq as soon as its own loop stops.
func Filter[V any](seq iter.Seq[V], keep func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if keep(v) && !yield(v) {
				return
			}
		}
	}
}

// Map yields f(v) for each value v of seq, and stops seq as soon as
// its own loop stops.
func Map[V, W any](seq iter.Seq[V], f func(V) W) iter.Seq[W] {
	return func(yield func(W) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Take yields the first n values of seq, and stops seq as soon as
// it has them or its own loop stops.
func Take[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		if n <= 0 {
			return
		}
		taken := 0
		for v := range seq {
			taken++
			if !yield(v) || taken == n {
				return
			}
		}
	}
}

// Zip yields the values of a and b in pairs, in order, until either
// ends, and stops both as soon as its own loop stops.
func Zip[V any](a, b iter.Seq[V]) iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		next, stop := iter.Pull(b)
		defer stop()
		for va := range a {
			vb, ok := next()
			if !ok || !yield(va, vb) {
				return
			}
		}
	}
}

// Merge yields the values of a and b, each sorted by cmp, as one
// sequence sorted by cmp, a's first on a tie, and stops both as
// soon as its own loop stops.
func Merge[V any](a, b iter.Seq[V], cmp func(V, V) int) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range a {
			if !yield(v) {
				return
			}
		}
		for v := range b {
			if !yield(v) {
				return
			}
		}
	}
}
