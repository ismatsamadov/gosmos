package main

import (
	"testing"
	"time"
)

func TestExhaustIn(t *testing.T) {
	day := 24 * time.Hour
	tests := []struct {
		name                  string
		remaining, lostLast7d int
		want                  time.Duration
		ok                    bool
	}{
		{"4 left, 2 lost this week", 4, 2, 14 * day, true},
		{"5 left, 1 lost this week", 5, 1, 35 * day, true},
		{"5 left, 3 lost this week: not whole weeks", 5, 3, 11*day + 16*time.Hour, true},
		{"nothing left", 0, 2, 0, true},
		{"overspent by two", -2, 1, 0, true},
		{"overspent, and nothing lost this week", -2, 0, 0, true},
		{"nothing lost this week", 5, 0, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := exhaustIn(tt.remaining, tt.lostLast7d)
			if got != tt.want || ok != tt.ok {
				t.Errorf("exhaustIn(%d, %d) = %v, %v, want %v, %v", tt.remaining, tt.lostLast7d, got, ok, tt.want, tt.ok)
			}
		})
	}
}
