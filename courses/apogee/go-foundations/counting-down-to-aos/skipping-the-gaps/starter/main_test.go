package main

import (
	"math"
	"testing"
)

func TestAverageBattery(t *testing.T) {
	tests := []struct {
		name     string
		from, to int
		want     float64
	}{
		{"the first two minutes", 0, 2, 7.91},
		{"the whole pass, two gaps", 0, 6, 7.87},
		{"a lone gap", 2, 3, 0},
		{"another lone gap", 4, 5, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := averageBattery(tt.from, tt.to)
			if math.IsNaN(got) || math.Abs(got-tt.want) > 0.01 {
				t.Errorf("averageBattery(%d, %d) = %.2f, want %.2f", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
