package main

import "testing"

func TestBatteryLine(t *testing.T) {
	tests := []struct {
		name  string
		battV float64
		want  string
	}{
		{"Pipit-1 after a sunlit orbit", 8.3, "FULL (8.30 V)"},
		{"Furrow-A right at FULL", 8.2, "FULL (8.20 V)"},
		{"Furrow-B a hair under FULL", 8.19, "OK (8.19 V)"},
		{"Cirrine-1 right at OK", 7.2, "OK (7.20 V)"},
		{"Brine-3 a hair under OK", 7.19, "LOW (7.19 V)"},
		{"Thistle right at LOW", 6.6, "LOW (6.60 V)"},
		{"Thistle a hair under LOW", 6.59, "CRITICAL (6.59 V)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := batteryLine(tt.battV)
			if got != tt.want {
				t.Errorf("batteryLine(%v) = %q, want %q", tt.battV, got, tt.want)
			}
		})
	}
}
