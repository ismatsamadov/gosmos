package main

import "testing"

func TestCheckPass(t *testing.T) {
	high := func(deg float64) bool { return deg >= 20 }
	low := func(deg float64) bool { return deg < 20 }
	for _, tc := range []struct {
		name string
		deg  float64
		rule func(float64) bool
		want bool
	}{
		{"below", 19, high, false}, {"boundary", 20, high, true},
		{"above", 60, high, true}, {"different rule", 5, low, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := checkPass(tc.deg, tc.rule); got != tc.want {
				t.Errorf("checkPass(%v, rule) = %v, want %v", tc.deg, got, tc.want)
			}
		})
	}
	calls := 0
	checkPass(41.5, func(deg float64) bool {
		calls++
		if deg != 41.5 {
			t.Errorf("rule argument = %v, want 41.5", deg)
		}
		return true
	})
	if calls != 1 {
		t.Errorf("rule calls = %d, want 1", calls)
	}
}
