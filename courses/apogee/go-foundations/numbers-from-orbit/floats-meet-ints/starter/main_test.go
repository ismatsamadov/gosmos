package main

import (
	"math"
	"testing"
)

func TestAverageMB(t *testing.T) {
	got := averageMB(1000, 3) // Furrowlight's three passes
	if math.Abs(got-333.333) > 0.001 {
		t.Errorf("averageMB(1000, 3) = %.3f, want 333.333", got)
	}
	got = averageMB(900, 3)
	if math.Abs(got-300) > 0.001 {
		t.Errorf("averageMB(900, 3) = %.3f, want 300.000", got)
	}
	got = averageMB(10, 4)
	if math.Abs(got-2.5) > 0.001 {
		t.Errorf("averageMB(10, 4) = %.3f, want 2.500", got)
	}
}
