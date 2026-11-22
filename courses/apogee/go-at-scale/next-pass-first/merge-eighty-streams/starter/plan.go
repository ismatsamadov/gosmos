package main

import (
	"cmp"
	"strings"
)

// byAOS orders passes by AOS, then by pass ID.
func byAOS(a, b PassView) int {
	return cmp.Or(a.AOS.Compare(b.AOS), strings.Compare(a.ID, b.ID))
}
