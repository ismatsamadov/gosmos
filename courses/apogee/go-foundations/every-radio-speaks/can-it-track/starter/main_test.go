package main

import "testing"

// Two radios main.go doesn't have, so the check must work for any radio, not only the three
// it knows by name.

// KaBandRadio is the Ka-band test radio on La Grande, which corrects Doppler shift itself.
type KaBandRadio struct{}

// RateMbps returns the Ka-band downlink rate.
func (k KaBandRadio) RateMbps() float64 { return 300 }

// TracksDoppler marks the Ka-band radio as correcting Doppler shift itself.
func (k KaBandRadio) TracksDoppler() {}

// SpareUHFRadio is the spare UHF radio in Port Ellery's stores, which can't correct Doppler
// shift either.
type SpareUHFRadio struct{}

// RateMbps returns the spare radio's downlink rate.
func (s SpareUHFRadio) RateMbps() float64 { return 0.0096 }

func TestNeedsPrecompensation(t *testing.T) {
	tests := []struct {
		name  string
		radio Radio
		want  bool
	}{
		{"S-band", SBandRadio{}, false},
		{"UHF", UHFRadio{}, true},
		{"X-band", XBandRadio{}, false},
		{"Ka-band", KaBandRadio{}, false},
		{"spare UHF", SpareUHFRadio{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := needsPrecompensation(tt.radio); got != tt.want {
				t.Errorf("needsPrecompensation(%s) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
