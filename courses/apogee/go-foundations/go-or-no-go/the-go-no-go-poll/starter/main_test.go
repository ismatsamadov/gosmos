package main

import "testing"

func TestPoll(t *testing.T) {
	tests := []struct {
		name       string
		up         bool
		windKmh    float64
		maxElevDeg float64
		battV      float64
		want       string
	}{
		{"all clear", true, 18, 41.5, 7.9, "GO"},
		{"station down", false, 18, 41.5, 7.9, "NO-GO: station down"},
		{"wind just under the stow limit", true, 71.9, 41.5, 7.9, "GO"},
		{"wind at the stow limit", true, 72, 41.5, 7.9, "NO-GO: wind"},
		{"peak under the mask", true, 18, 2.5, 7.9, "NO-GO: below mask"},
		{"battery under LOW", true, 18, 41.5, 6.5, "NO-GO: battery"},
		{"peak exactly at the mask", true, 18, 3, 7.9, "GO"},
		{"battery exactly at LOW", true, 18, 41.5, 6.6, "GO"},
		{"down and windy", false, 90, 41.5, 7.9, "NO-GO: station down"},
		{"windy, low and weak", true, 80, 2.5, 6.4, "NO-GO: wind"},
		{"low pass, weak battery", true, 18, 2.5, 6.4, "NO-GO: below mask"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := poll(tt.up, tt.windKmh, tt.maxElevDeg, tt.battV)
			if got != tt.want {
				t.Errorf("poll(%v, %v, %v, %v) = %q, want %q", tt.up, tt.windKmh, tt.maxElevDeg, tt.battV, got, tt.want)
			}
		})
	}
}
