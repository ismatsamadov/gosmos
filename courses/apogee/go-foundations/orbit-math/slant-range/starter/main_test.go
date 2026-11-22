package main

import (
	"math"
	"testing"
)

func TestSlantRangeKm(t *testing.T) {
	tests := []struct {
		name           string
		altKm, elevDeg float64
		want           float64
	}{
		{"Pipit-1 overhead", 525, 90, 525.00},
		{"Pipit-1 at the horizon", 525, 0, 2639.16},
		{"Pipit-1 at 30 deg", 525, 30, 951.24},
		{"Thistle at 10 deg", 415, 10, 1479.05},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slantRangeKm(tt.altKm, tt.elevDeg)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("slantRangeKm(%v, %v) = %.2f, want %.2f", tt.altKm, tt.elevDeg, got, tt.want)
			}
		})
	}
}
