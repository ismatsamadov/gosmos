package main

import "slices"

// byAOS orders passes by AOS.
func byAOS(a, b PassView) int {
	return a.AOS.Compare(b.AOS)
}

// canServe reports whether the station offers one of the satellite's bands.
func canServe(s Site, sat Satellite) bool {
	return slices.ContainsFunc(sat.Bands, func(b Band) bool { return slices.Contains(s.Bands, b) })
}
