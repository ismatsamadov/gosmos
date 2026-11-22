package main

import (
	"math"
	"testing"
)

func TestDopplerHz(t *testing.T) {
	tests := []struct {
		name                 string
		freqHz, rangeRateMps float64
		want                 float64
	}{
		{"Pipit-1 rising", 437.25e6, -7000, 10209.56},
		{"Pipit-1 setting", 437.25e6, 7000, -10209.56},
		{"Pipit-2 rising", 2245e6, -6500, 48675.34},
		{"Pipit-2 at closest approach", 2245e6, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dopplerHz(tt.freqHz, tt.rangeRateMps)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("dopplerHz(%.0f, %.0f) = %.2f, want %.2f", tt.freqHz, tt.rangeRateMps, got, tt.want)
			}
		})
	}
}
