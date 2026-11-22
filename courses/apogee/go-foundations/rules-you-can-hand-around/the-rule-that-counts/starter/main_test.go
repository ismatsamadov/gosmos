package main

import "testing"

func TestMakeCounter(t *testing.T) {
	a, b := makeCounter(), makeCounter()
	if a == nil || b == nil {
		t.Fatal("makeCounter() = nil, want a function")
	}
	for _, tc := range []struct {
		counter func() int
		want    int
	}{{a, 1}, {a, 2}, {b, 1}, {a, 3}, {b, 2}} {
		if got := tc.counter(); got != tc.want {
			t.Errorf("counter() = %d, want %d", got, tc.want)
		}
	}
}
