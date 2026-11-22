package main

import (
	"math"
	"testing"
)

func TestPathLossDB(t *testing.T) {
	tests := []struct {
		name                    string
		altKm, elevDeg, freqMHz float64
		want                    float64
	}{
		{"Pipit-1 at the horizon", 525, 0, 437.25, 153.68},
		{"Pipit-1 overhead", 525, 90, 437.25, 139.66},
		{"Pipit-2 at 10 deg", 550, 10, 2245, 164.64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pathLossDB(tt.altKm, tt.elevDeg, tt.freqMHz)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("pathLossDB(%v, %v, %v) = %.2f, want %.2f", tt.altKm, tt.elevDeg, tt.freqMHz, got, tt.want)
			}
		})
	}
}
