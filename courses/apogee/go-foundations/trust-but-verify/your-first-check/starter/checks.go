package main

import "testing"

// checkBand checks a version of batteryBand on voltages around each threshold.
func checkBand(t *testing.T, band func(v float64) string) {
	// TODO: call band at 8.2, 7.2 and 6.6 V and 0.01 V below each; t.Errorf on each wrong band.
}
