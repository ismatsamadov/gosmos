package main

import (
	"testing"
	"time"
)

func TestPassDuration(t *testing.T) {
	for _, tc := range []struct {
		name, aos, los string
		want           time.Duration
	}{
		{"across midnight", "2027-09-22T23:54:00Z", "2027-09-23T00:05:00Z", 11 * time.Minute},
		{"van clock", "2027-09-22T16:54:00-07:00", "2027-09-22T17:05:00-07:00", 11 * time.Minute},
		{"mixed clocks", "2027-09-22T16:54:00-07:00", "2027-09-23T00:05:00Z", 11 * time.Minute},
		{"a pass to the second", "2027-09-23T04:12:00Z", "2027-09-23T04:21:12Z", 9*time.Minute + 12*time.Second},
		{"AOS is LOS", "2027-09-23T00:05:00Z", "2027-09-23T00:05:00Z", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := passDuration(tc.aos, tc.los)
			if got != tc.want || err != nil {
				t.Errorf("passDuration(%q, %q) = %v, %v, want %v, nil", tc.aos, tc.los, got, err, tc.want)
			}
		})
	}
	for _, tc := range []struct{ name, aos, los string }{
		{"LOS before AOS", "2027-09-23T00:05:00Z", "2027-09-22T23:54:00Z"},
		{"clock readings only", "23:54", "00:05"},
		{"bad AOS", "2027-09-22 23:54", "2027-09-23T00:05:00Z"},
		{"bad LOS", "2027-09-22T23:54:00Z", "2027-09-23 00:05"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := passDuration(tc.aos, tc.los)
			if got != 0 || err == nil {
				t.Errorf("passDuration(%q, %q) = %v, %v, want 0 and an error", tc.aos, tc.los, got, err)
			}
		})
	}
}
