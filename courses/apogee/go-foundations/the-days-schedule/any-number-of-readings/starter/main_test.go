package main

import (
	"math"
	"testing"
)

func TestAverage(t *testing.T) {
	temps := []float64{-4.5, -3.9, -2.4}
	tests := []struct {
		call      string
		got, want float64
	}{
		{"average()", average(), 0},
		{"average(7.9)", average(7.9), 7.9},
		{"average(1, 2, 3, 4)", average(1, 2, 3, 4), 2.5},
		{"average(temps...)", average(temps...), -3.6},
	}
	for _, tt := range tests {
		t.Run(tt.call, func(t *testing.T) {
			if math.IsNaN(tt.got) || math.Abs(tt.got-tt.want) > 0.01 {
				t.Errorf("%s = %.2f, want %.2f", tt.call, tt.got, tt.want)
			}
		})
	}
}
