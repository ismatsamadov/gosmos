package main

import "testing"

func TestIsPolar(t *testing.T) {
	tests := []struct {
		name string
		sat  Satellite
		want bool
	}{
		{"Pipit-1", Satellite{ID: "pipit-1", Name: "Pipit-1", Orbit: Orbit{AltKm: 525, IncDeg: 97.5}}, true},
		{"Brine-3", Satellite{ID: "brine-3", Name: "Brine-3", Orbit: Orbit{AltKm: 480, IncDeg: 53}}, false},
		{"Thistle", Satellite{ID: "thistle", Name: "Thistle", Orbit: Orbit{AltKm: 415, IncDeg: 51.6}}, false},
		{"exactly 80", Satellite{ID: "test-80", Name: "Test-80", Orbit: Orbit{AltKm: 500, IncDeg: 80}}, true},
		{"exactly 100", Satellite{ID: "test-100", Name: "Test-100", Orbit: Orbit{AltKm: 500, IncDeg: 100}}, true},
		{"just under 80", Satellite{ID: "test-79", Name: "Test-79", Orbit: Orbit{AltKm: 500, IncDeg: 79.9}}, false},
		{"just over 100", Satellite{ID: "test-101", Name: "Test-101", Orbit: Orbit{AltKm: 500, IncDeg: 100.1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPolar(tt.sat); got != tt.want {
				t.Errorf("isPolar(%+v) = %v, want %v", tt.sat, got, tt.want)
			}
		})
	}
}
