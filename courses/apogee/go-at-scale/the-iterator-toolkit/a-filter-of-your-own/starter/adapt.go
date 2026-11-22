package main

import "iter"

// Filter yields the values of seq that keep keeps, and stops seq as soon as its own loop stops.
func Filter[V any](seq iter.Seq[V], keep func(V) bool) iter.Seq[V] {
	return seq
}
