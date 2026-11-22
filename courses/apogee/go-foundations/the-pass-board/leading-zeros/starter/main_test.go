package main

import "testing"

func TestClock(t *testing.T) {
	tests := []struct {
		name string
		h, m int
		want string
	}{
		{"Brine-3 AOS at Lantern Atoll", 4, 7, "04:07"},
		{"last minute of the day", 23, 59, "23:59"},
		{"midnight", 0, 0, "00:00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clock(tt.h, tt.m); got != tt.want {
				t.Errorf("clock(%d, %d) = %q, want %q", tt.h, tt.m, got, tt.want)
			}
		})
	}
}
