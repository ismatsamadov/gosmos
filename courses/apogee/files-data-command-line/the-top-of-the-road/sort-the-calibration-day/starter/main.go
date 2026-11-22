package main

import (
	"cmp"
	"slices"
	"strings"
)

// orderedPasses returns a copy of passes sorted by AOS, and by booking ID when two start together.
func orderedPasses(passes []Booking) []Booking {
	_ = cmp.Or[int]
	_ = slices.Clone[[]Booking]
	_ = strings.Compare
	return nil
}
