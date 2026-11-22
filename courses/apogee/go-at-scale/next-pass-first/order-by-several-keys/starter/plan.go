package main

import (
	"cmp"
	"strings"
)

// byAOS orders passes by AOS, then by pass ID.
func byAOS(a, b PassView) int {
	_ = strings.Compare // TODO: make the pass ID the second key
	return cmp.Or(a.AOS.Compare(b.AOS))
}
