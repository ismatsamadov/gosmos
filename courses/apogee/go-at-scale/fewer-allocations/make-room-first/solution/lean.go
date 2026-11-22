package main

import "iter"

// collect gathers seq's passes into a slice, in order. n is how many
// passes the station had on the last run.
func collect(seq iter.Seq[PassView], n int) []PassView {
	out := make([]PassView, 0, n)
	for p := range seq {
		out = append(out, p)
	}
	return out
}
