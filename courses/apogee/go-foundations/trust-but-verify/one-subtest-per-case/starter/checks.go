package main

import "testing"

// windCase is one wind speed and the call Cape Wren's dish must make in it.
type windCase struct {
	name string
	kmh  float64
	want string
}

// windCases are the cases checkWind runs.
var windCases = []windCase{
	{"calm-30", 30, "track"},
	{"warning-50", 50, "track with warning"},
	{"warning-71.9", 71.9, "track with warning"},
	{"stow-72", 72, "stow"},
	{"storm-95", 95, "stow"},
}

// checkWind checks a version of windCall on every case in windCases.
func checkWind(t *testing.T, windCall func(kmh float64) string) {
	for _, c := range windCases {
		if got := windCall(c.kmh); got != c.want {
			t.Errorf("windCall(%v) = %q, want %q", c.kmh, got, c.want)
		}
	}
}
