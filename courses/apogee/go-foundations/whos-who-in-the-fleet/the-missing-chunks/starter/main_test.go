package main

import (
	"slices"
	"testing"
)

func TestMissing(t *testing.T) {
	cases := []struct {
		name   string
		chunks map[int]string
		total  int
		want   []int
	}{
		{"furrow-a image, two lost", map[int]string{0: "img-0", 1: "img-1", 3: "img-3"}, 5, []int{2, 4}},
		{"everything came down", map[int]string{0: "img-0", 1: "img-1", 2: "img-2"}, 3, nil},
		{"nothing came down", map[int]string{}, 3, []int{0, 1, 2}},
		{"only the last one lost", map[int]string{0: "img-0", 1: "img-1", 2: "img-2", 3: "img-3"}, 5, []int{4}},
		{"chunk 1 came down empty", map[int]string{0: "img-0", 1: "", 2: "img-2"}, 3, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := missing(c.chunks, c.total); !slices.Equal(got, c.want) {
				t.Errorf("missing(%v, %d) = %v, want %v", c.chunks, c.total, got, c.want)
			}
		})
	}
}
