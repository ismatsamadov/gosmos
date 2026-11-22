package main

import "testing"

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"Hollis at Furrowlight", "hollis.vey@furrowlight.example", "h***@furrowlight.example"},
		{"Brinewatch's ops desk", "ops@brinewatch.example", "o***@brinewatch.example"},
		{"Cirrine's uplink team", "uplink@cirrine.example", "u***@cirrine.example"},
		{"the Fishbowl's own", "fishbowl@apogee.example", "f***@apogee.example"},
		{"a one-letter name", "h@furrowlight.example", "h***@furrowlight.example"},
		{"no @ at all", "hollis.vey.furrowlight.example", "[redacted]"},
		{"nothing before the @", "@furrowlight.example", "[redacted]"},
		{"nothing after the @", "hollis.vey@", "[redacted]"},
		{"two @", "hollis@vey@furrowlight.example", "[redacted]"},
		{"an empty field", "", "[redacted]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maskEmail(tt.s); got != tt.want {
				t.Errorf("maskEmail(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}
