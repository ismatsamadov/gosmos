package main

import (
	"slices"
	"testing"
)

func TestDedupe(t *testing.T) {
	cases := []struct {
		name string
		seqs []int
		want []int
	}{
		{"thistle twice over", []int{4182, 4183, 4182, 4184, 4183}, []int{4182, 4183, 4184}},
		{"first appearance wins", []int{4184, 4182, 4184}, []int{4184, 4182}},
		{"no repeats", []int{7455, 7456, 7457}, []int{7455, 7456, 7457}},
		{"silent pass", []int{}, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// A map's order changes from run to run, so one lucky run proves nothing.
			for range 20 {
				if got := dedupe(slices.Clone(c.seqs)); !slices.Equal(got, c.want) {
					t.Fatalf("dedupe(%v) = %v, want %v", c.seqs, got, c.want)
				}
			}
		})
	}
}
