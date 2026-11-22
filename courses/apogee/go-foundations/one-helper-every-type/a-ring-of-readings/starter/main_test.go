package main

import (
	"slices"
	"testing"
)

func TestRingInts(t *testing.T) {
	cases := []struct {
		name   string
		pushes int // pushes 1, 2, ... up to this number
		want   []int
	}{
		{"two-of-three", 2, []int{1, 2}},
		{"exactly-full", 3, []int{1, 2, 3}},
		{"five-into-three", 5, []int{3, 4, 5}},
		{"seven-into-three", 7, []int{5, 6, 7}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := NewRing[int](3)
			for i := 1; i <= c.pushes; i++ {
				r.Push(i)
			}
			if got := r.Items(); !slices.Equal(got, c.want) {
				t.Errorf("ring of 3 after pushing 1 to %d: Items() = %v, want %v", c.pushes, got, c.want)
			}
		})
	}
}

func TestRingPassIDs(t *testing.T) {
	r := NewRing[string](2)
	for _, id := range []string{"pipit-1/PEL/04182", "thistle/LAN/04183", "furrow-a/SQT/04184"} {
		r.Push(id)
	}
	want := []string{"thistle/LAN/04183", "furrow-a/SQT/04184"}
	if got := r.Items(); !slices.Equal(got, want) {
		t.Errorf("ring of 2 after three pass IDs: Items() = %q, want %q", got, want)
	}
}

// TestItemsIsACopy changes the slice Items returns: the ring must not change with it.
func TestItemsIsACopy(t *testing.T) {
	r := NewRing[int](3)
	for i := 1; i <= 3; i++ {
		r.Push(i)
	}
	items := r.Items()
	if len(items) != 3 {
		t.Fatalf("ring of 3 after pushing 1 to 3: Items() = %v, want 3 items", items)
	}
	items[0] = 99
	if got, want := r.Items(), []int{1, 2, 3}; !slices.Equal(got, want) {
		t.Errorf("after setting Items()[0] = 99, Items() = %v, want %v", got, want)
	}
}
