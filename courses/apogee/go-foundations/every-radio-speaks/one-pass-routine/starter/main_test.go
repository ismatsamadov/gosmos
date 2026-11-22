package main

import (
	"math"
	"testing"
)

func TestDownlinkSeconds(t *testing.T) {
	tests := []struct {
		name      string
		radio     Radio
		megabytes float64
		want      float64
	}{
		{"X-band", XBandRadio{}, 9000, 480},
		{"S-band", SBandRadio{}, 30, 60},
		{"UHF", UHFRadio{}, 0.0072, 6},
	}
	for _, tt := range tests {
		got := downlinkSeconds(tt.radio, tt.megabytes)
		if math.Abs(got-tt.want) > 0.01 {
			t.Errorf("downlinkSeconds(%s, %v) = %.2f, want %.2f", tt.name, tt.megabytes, got, tt.want)
		}
	}
}
