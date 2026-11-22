package main

import "testing"

func TestElevationCall(t *testing.T) {
	tests := []struct {
		name       string
		maxElevDeg float64
		maskDeg    float64
		want       string
	}{
		{"Furrow-A high over Ystmark", 42, 5, "GO (max 42.0 deg, mask 5.0 deg)"},
		{"Thistle low over Cape Wren", 4.5, 5, "NO-GO (max 4.5 deg, mask 5.0 deg)"},
		{"Pipit-2 peaking at the mask", 5, 5, "GO (max 5.0 deg, mask 5.0 deg)"},
		{"Brine-3 low over Salar Quieto", 4.5, 3, "GO (max 4.5 deg, mask 3.0 deg)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := elevationCall(tt.maxElevDeg, tt.maskDeg)
			if got != tt.want {
				t.Errorf("elevationCall(%v, %v) = %q, want %q", tt.maxElevDeg, tt.maskDeg, got, tt.want)
			}
		})
	}
}
