package main

import (
	"cmp"
	"slices"
	"strings"
)

// byAOS orders passes by AOS, then by pass ID.
func byAOS(a, b PassView) int {
	return cmp.Or(a.AOS.Compare(b.AOS), strings.Compare(a.ID, b.ID))
}

// canServe reports whether the station offers one of the satellite's bands.
func canServe(s Site, sat Satellite) bool {
	return slices.ContainsFunc(sat.Bands, func(b Band) bool { return slices.Contains(s.Bands, b) })
}
