package main

import "testing"

func TestStationLine(t *testing.T) {
	tests := []struct {
		name string
		s    Station
		want string
	}{
		{"Port Ellery", Station{Code: "PEL", Name: "Port Ellery", LatDeg: 46.9, MaskDeg: 5}, "PEL Port Ellery (46.9 N, mask 5.0)"},
		{"Salar Quieto", Station{Code: "SQT", Name: "Salar Quieto", LatDeg: -23.4, MaskDeg: 3}, "SQT Salar Quieto (23.4 S, mask 3.0)"},
		{"Lantern Atoll", Station{Code: "LAN", Name: "Lantern Atoll", LatDeg: -1.8, MaskDeg: 8}, "LAN Lantern Atoll (1.8 S, mask 8.0)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stationLine(tt.s); got != tt.want {
				t.Errorf("stationLine(%v) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}
