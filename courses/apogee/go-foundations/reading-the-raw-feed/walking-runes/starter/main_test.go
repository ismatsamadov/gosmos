package main

import "testing"

func TestCallSign(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{"station name", "Salar Quieto", "SALAR-QUIETO"},
		{"one rune, two bytes", "Añil", "A_IL"},
		{"digits and hyphens", "Pipit-3", "PIPIT-3"},
		{"ends of each range", "az AZ 09", "AZ-AZ-09"},
		{"empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := callSign(c.in); got != c.want {
				t.Errorf("callSign(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
