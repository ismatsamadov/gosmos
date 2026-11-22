package main

import "iter"

// Change is one difference between two schedules.
type Change struct {
	Kind    string  // "old only", "new only" or "dish"
	Contact Contact // the old planner's for "old only", the new planner's otherwise
}

// Compare yields every difference between the old planner's
// contacts and the new planner's, both in byAOS order. It walks
// them a second at a time and yields what changes finds among
// each second's contacts.
func Compare(old, new iter.Seq[Contact]) iter.Seq[Change] {
	return func(yield func(Change) bool) {}
}
