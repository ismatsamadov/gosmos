package main

import (
	"slices"
	"testing"
)

func TestReport(t *testing.T) {
	cases := []struct {
		name   string
		counts map[string]int
		want   []string
	}{
		{
			name:   "busy hour",
			counts: map[string]int{"thistle": 2, "pipit-1": 3, "brine-3": 1},
			want:   []string{"brine-3: 1", "pipit-1: 3", "thistle: 2"},
		},
		{
			name: "whole fleet",
			counts: map[string]int{
				"pipit-3": 1, "thistle": 2, "furrow-b": 5, "pipit-1": 3,
				"cirrine-1": 4, "brine-3": 1, "pipit-2": 4, "furrow-a": 6,
			},
			want: []string{
				"brine-3: 1", "cirrine-1: 4", "furrow-a: 6", "furrow-b: 5",
				"pipit-1: 3", "pipit-2: 4", "pipit-3: 1", "thistle: 2",
			},
		},
		{
			name:   "quiet hour",
			counts: map[string]int{},
			want:   nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// A map's order changes from run to run, so one lucky run proves nothing.
			for range 20 {
				if got := report(c.counts); !slices.Equal(got, c.want) {
					t.Fatalf("report(%v) = %q, want %q", c.counts, got, c.want)
				}
			}
		})
	}
}
