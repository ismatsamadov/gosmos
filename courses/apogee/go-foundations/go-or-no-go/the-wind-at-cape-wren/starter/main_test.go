package main

import "testing"

func TestWindCall(t *testing.T) {
	tests := []struct {
		name    string
		windKmh float64
		want    string
	}{
		{"a calm morning", 30, "track"},
		{"wind just under the warning", 49.9, "track"},
		{"the warning starts", 50, "track with warning"},
		{"a gust under the stow limit", 71.9, "track with warning"},
		{"the stow starts", 72, "stow"},
		{"a southern gale", 95, "stow"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windCall(tt.windKmh)
			if got != tt.want {
				t.Errorf("windCall(%v) = %q, want %q", tt.windKmh, got, tt.want)
			}
		})
	}
}
