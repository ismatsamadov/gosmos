package main

import "testing"

// checkBand checks a version of batteryBand on voltages around each threshold.
func checkBand(t *testing.T, band func(v float64) string) {
	cases := []struct {
		v    float64
		want string
	}{
		{8.2, "FULL"}, {8.19, "OK"},
		{7.2, "OK"}, {7.19, "LOW"},
		{6.6, "LOW"}, {6.59, "CRITICAL"},
	}
	for _, c := range cases {
		if got := band(c.v); got != c.want {
			t.Errorf("band(%.2f) = %q, want %q", c.v, got, c.want)
		}
	}
}
