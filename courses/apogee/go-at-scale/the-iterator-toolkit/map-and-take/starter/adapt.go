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
	return func(yield func(W) bool) {}
}

// Take yields the first n values of seq, and stops seq as soon as
// it has them or its own loop stops.
func Take[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {}
}
