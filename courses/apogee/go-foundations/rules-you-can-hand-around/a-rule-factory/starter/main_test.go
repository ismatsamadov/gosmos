package main

import "testing"

func TestMinElevation(t *testing.T) {
	furrow, cirrine := minElevation(20), minElevation(30)
	if furrow == nil || cirrine == nil {
		t.Fatal("minElevation() = nil, want a function")
	}
	for _, tc := range []struct {
		name string
		rule func(float64) bool
		deg  float64
		want bool
	}{
		{"furrow below", furrow, 19.9, false}, {"furrow boundary", furrow, 20, true},
		{"cirrine below", cirrine, 25, false}, {"cirrine boundary", cirrine, 30, true},
		{"furrow retained", furrow, 25, true}, {"zero", minElevation(0), 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.rule(tc.deg); got != tc.want {
				t.Errorf("rule(%v) = %v, want %v", tc.deg, got, tc.want)
			}
		})
	}
}
