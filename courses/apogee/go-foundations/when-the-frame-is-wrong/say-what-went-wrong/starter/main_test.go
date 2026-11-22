package main

import "testing"

func TestCheckBattery(t *testing.T) {
	tests := []struct {
		sat  string
		v    float64
		want string // "" means nil
	}{
		{"thistle", 6.2, "thistle: battery 6.20 V below 6.60 V"},
		{"brine-3", 6.55, "brine-3: battery 6.55 V below 6.60 V"},
		{"pipit-1", 7.9, ""},
		{"pipit-2", 6.6, ""},
	}
	for _, tt := range tests {
		t.Run(tt.sat, func(t *testing.T) {
			err := checkBattery(tt.sat, tt.v)
			switch {
			case tt.want == "" && err != nil:
				t.Errorf("checkBattery(%q, %v) = %q, want nil", tt.sat, tt.v, err)
			case tt.want != "" && err == nil:
				t.Errorf("checkBattery(%q, %v) = nil, want %q", tt.sat, tt.v, tt.want)
			case tt.want != "" && err.Error() != tt.want:
				t.Errorf("checkBattery(%q, %v) = %q, want %q", tt.sat, tt.v, err, tt.want)
			}
		})
	}
}
