package main

import (
	"cmp"
	"slices"
	"strings"
)

// orderedPasses returns a copy of passes sorted by AOS, and by booking ID when two start together.
func orderedPasses(passes []Booking) []Booking {
	sorted := slices.Clone(passes)
	slices.SortFunc(sorted, func(a, b Booking) int {
		return cmp.Or(a.AOS.Compare(b.AOS), strings.Compare(a.ID, b.ID))
	})
	return sorted
}
