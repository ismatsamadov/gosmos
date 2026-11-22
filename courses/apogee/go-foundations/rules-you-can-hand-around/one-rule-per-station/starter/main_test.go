package main

import "testing"

func TestStationRules(t *testing.T) {
	first, second, third := stationRules()
	for i, rule := range []func(int) bool{first, second, third} {
		if rule == nil {
			t.Fatalf("stationRules() rule %d = nil, want a function", i+1)
		}
		for deg := 0; deg <= 4; deg++ {
			want := deg >= i+1
			if got := rule(deg); got != want {
				t.Errorf("rule %d(%d) = %v, want %v", i+1, deg, got, want)
			}
		}
	}
}
