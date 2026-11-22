package main

import "testing"

func TestAOSMinute(t *testing.T) {
	tests := []struct {
		name    string
		maskDeg float64
		want    int
	}{
		{"Ystmark's 5 degree mask", 5, 1},
		{"a 30 degree mask", 30, 3},
		{"a mask at the peak", 64, 7},
		{"a mask above the peak", 70, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aosMinute(tt.maskDeg)
			if got != tt.want {
				t.Errorf("aosMinute(%v) = %v, want %v", tt.maskDeg, got, tt.want)
			}
		})
	}
}
