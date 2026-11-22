package main

import "testing"

// pipit1Beacon is a beacon from Pipit-1: sequence 4182, battery 7910 mV.
const pipit1Beacon = "4150011056001ee6ffd30082"

func TestDecodeBeacon(t *testing.T) {
	seq, battMV, err := decodeBeacon(pipit1Beacon)
	if err != nil {
		t.Fatalf("decodeBeacon(%q) error = %v, want nil", pipit1Beacon, err)
	}
	if seq != 4182 || battMV != 7910 {
		t.Errorf("decodeBeacon(%q) = %d, %d, want 4182, 7910", pipit1Beacon, seq, battMV)
	}
}

func TestDecodeBeaconRejects(t *testing.T) {
	cases := []struct{ name, hexText string }{
		{"not hex", "4150011056001ee6ffd3008g"},
		{"empty", ""},
		{"8 bytes", "4150011056001ee6"},
		{"13 bytes", "4150011056001ee6ffd3008200"},
		{"sync bytes AQ", "4151011056001ee6ffd30083"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, _, err := decodeBeacon(c.hexText); err == nil {
				t.Errorf("decodeBeacon(%q) error = nil, want an error", c.hexText)
			}
		})
	}
}
