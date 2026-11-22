package main

import "testing"

func TestOperatorOf(t *testing.T) {
	reg := map[string]string{
		"pipit-1": "Apogee",
		"pipit-3": "",
		"thistle": "Kessick Polytechnic",
	}
	cases := []struct {
		id, want string
	}{
		{"pipit-1", "Apogee"},
		{"thistle", "Kessick Polytechnic"},
		{"brine-2", "unknown satellite"},
		{"pipit-3", ""},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			if got := operatorOf(reg, c.id); got != c.want {
				t.Errorf("operatorOf(reg, %q) = %q, want %q", c.id, got, c.want)
			}
		})
	}
}
