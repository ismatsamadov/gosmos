package main

import (
	"math"
	"testing"
)

func TestPeriodMinutes(t *testing.T) {
	got := periodMinutes(525) // Pipit-1
	if math.Abs(got-94.985) > 0.01 {
		t.Errorf("periodMinutes(525) = %.3f, want 94.985", got)
	}
	got = periodMinutes(415) // Thistle
	if math.Abs(got-92.722) > 0.01 {
		t.Errorf("periodMinutes(415) = %.3f, want 92.722", got)
	}
	got = periodMinutes(600) // Cirrine-1
	if math.Abs(got-96.539) > 0.01 {
		t.Errorf("periodMinutes(600) = %.3f, want 96.539", got)
	}
}
