package main

import "testing"

func TestSeverityOrder(t *testing.T) {
	if !(Info < Warning && Warning < Critical) {
		t.Errorf("Info, Warning, Critical = %d, %d, %d, want them in increasing order", Info, Warning, Critical)
	}
}

func TestWorst(t *testing.T) {
	cases := []struct {
		name       string
		a, b, want Severity
	}{
		{"late frame, then critical battery", Info, Critical, Critical},
		{"low battery, then late frame", Warning, Info, Warning},
		{"overheating, then low battery", Critical, Warning, Critical},
		{"two low batteries", Warning, Warning, Warning},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := worst(c.a, c.b); got != c.want {
				t.Errorf("worst(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
			}
		})
	}
}
