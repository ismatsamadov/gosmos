package main

import "testing"

func TestCheckBattery(t *testing.T) {
	tests := []struct {
		v    float64
		want string // "" means nil
	}{
		{7.4, ""},
		{6.6, ""},
		{6.2, "battery too low"},
	}
	for _, tt := range tests {
		err := checkBattery(tt.v)
		switch {
		case tt.want == "" && err != nil:
			t.Errorf("checkBattery(%v) = %q, want nil", tt.v, err)
		case tt.want != "" && err == nil:
			t.Errorf("checkBattery(%v) = nil, want %q", tt.v, tt.want)
		case tt.want != "" && err.Error() != tt.want:
			t.Errorf("checkBattery(%v) = %q, want %q", tt.v, err, tt.want)
		}
	}
}
