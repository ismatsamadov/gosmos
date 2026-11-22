package main

import (
	"math"
	"testing"
)

func TestAvailability(t *testing.T) {
	tests := []struct {
		name          string
		flown, booked int
		want          float64
	}{
		{"Lantern Atoll's quarter", 1177, 1199, 0.98165},
		{"Brinewatch's first exact week", 61, 62, 0.98387},
		{"a day with every pass flown", 15, 15, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := availability(tt.flown, tt.booked)
			if err != nil {
				t.Fatalf("availability(%d, %d) error = %q, want none", tt.flown, tt.booked, err)
			}
			if math.Abs(got-tt.want) > 0.000005 {
				t.Errorf("availability(%d, %d) = %.5f, want %.5f", tt.flown, tt.booked, got, tt.want)
			}
		})
	}
}

func TestAvailabilityRefuses(t *testing.T) {
	tests := []struct {
		name          string
		flown, booked int
		want          string
	}{
		{"nothing booked", 0, 0, "no booked passes"},
		{"more flown than booked", 5, 4, "more flown than booked"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := availability(tt.flown, tt.booked)
			if err == nil {
				t.Fatalf("availability(%d, %d) = %v, <nil>, want error %q", tt.flown, tt.booked, got, tt.want)
			}
			if err.Error() != tt.want {
				t.Errorf("availability(%d, %d) error = %q, want %q", tt.flown, tt.booked, err, tt.want)
			}
			if got != 0 {
				t.Errorf("availability(%d, %d) = %v with its error, want 0", tt.flown, tt.booked, got)
			}
		})
	}
}
