package main

import "testing"

func TestHoursMinutes(t *testing.T) {
	tests := []struct {
		name         string
		total        int
		wantH, wantM int
	}{
		{"Pipit-1 AOS at Port Ellery", 252, 4, 12},
		{"midnight", 0, 0, 0},
		{"last minute of the day", 1439, 23, 59},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, m := hoursMinutes(tt.total)
			if h != tt.wantH || m != tt.wantM {
				t.Errorf("hoursMinutes(%d) = %d, %d, want %d, %d", tt.total, h, m, tt.wantH, tt.wantM)
			}
		})
	}
}
