package main

import (
	"math"
	"testing"
)

func TestBurnRate(t *testing.T) {
	tests := []struct {
		name         string
		target       Target
		booked, lost int
		want         float64
	}{
		{"two of Lantern Atoll's four passes in six hours", Atoll, 4, 2, 50},
		{"the night of 27-28 December: six of eight", Atoll, 8, 6, 75},
		{"one pass in 200 at 99.5%", Standard, 200, 1, 1},
		{"nothing booked", Atoll, 0, 0, 0},
		{"a loss against 10000", 10000, 4, 1, math.Inf(1)},
		{"no loss against 10000", 10000, 4, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := burnRate(tt.target, tt.booked, tt.lost)
			ok := math.Abs(got-tt.want) <= 0.01
			if math.IsInf(tt.want, 1) {
				ok = math.IsInf(got, 1)
			}
			if !ok {
				t.Errorf("burnRate(%d, %d, %d) = %.2f, want %.2f", tt.target, tt.booked, tt.lost, got, tt.want)
			}
		})
	}
}
