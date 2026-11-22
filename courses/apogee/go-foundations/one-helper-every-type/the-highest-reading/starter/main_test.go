package main

import (
	"math"
	"testing"
)

func TestMaxInts(t *testing.T) {
	got, ok := Max([]int{3, 9, 4})
	if got != 9 || !ok {
		t.Errorf("Max([3 9 4]) = %v, %v, want 9, true", got, ok)
	}
}

func TestMaxFloats(t *testing.T) {
	got, ok := Max([]float64{61.4, 23, 44.2})
	if math.Abs(got-61.4) > 0.01 || !ok {
		t.Errorf("Max([61.4 23 44.2]) = %.2f, %v, want 61.40, true", got, ok)
	}
}

// TestMaxBelowZero uses Ystmark's battery temperatures in the polar night.
func TestMaxBelowZero(t *testing.T) {
	got, ok := Max([]float64{-4.5, -12.25, -7})
	if math.Abs(got-(-4.5)) > 0.01 || !ok {
		t.Errorf("Max([-4.5 -12.25 -7]) = %.2f, %v, want -4.50, true", got, ok)
	}
}

func TestMaxStrings(t *testing.T) {
	got, ok := Max([]string{"YST", "CWR", "SQT"})
	if got != "YST" || !ok {
		t.Errorf("Max([YST CWR SQT]) = %q, %v, want %q, true", got, ok, "YST")
	}
}

func TestMaxEmpty(t *testing.T) {
	got, ok := Max([]int{})
	if got != 0 || ok {
		t.Errorf("Max([]) = %v, %v, want 0, false", got, ok)
	}
}
