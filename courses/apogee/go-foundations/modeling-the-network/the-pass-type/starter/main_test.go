package main

import (
	"testing"
	"time"
)

// at returns the given April 2027 day and time of day, in UTC.
func at(day, h, m, s int) time.Time {
	return time.Date(2027, time.April, day, h, m, s, 0, time.UTC)
}

func TestDuration(t *testing.T) {
	tests := []struct {
		name string
		p    Pass
		want time.Duration
	}{
		{"Pipit-1 over Port Ellery", Pass{SatID: "pipit-1", Station: "PEL", AOS: at(15, 4, 12, 0), LOS: at(15, 4, 21, 12)}, 9*time.Minute + 12*time.Second},
		{"AOS equal to LOS", Pass{SatID: "thistle", Station: "LAN", AOS: at(15, 6, 47, 0), LOS: at(15, 6, 47, 0)}, 0},
		{"Thistle across midnight", Pass{SatID: "thistle", Station: "LAN", AOS: at(15, 23, 58, 0), LOS: at(16, 0, 7, 0)}, 9 * time.Minute},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := duration(tt.p); got != tt.want {
				t.Errorf("duration(AOS %s, LOS %s) = %v, want %v", tt.p.AOS.Format(time.RFC3339), tt.p.LOS.Format(time.RFC3339), got, tt.want)
			}
		})
	}
}
