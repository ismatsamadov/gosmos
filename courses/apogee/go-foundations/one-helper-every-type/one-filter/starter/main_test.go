package main

import (
	"slices"
	"testing"
)

func TestFilterPasses(t *testing.T) {
	passes := []Pass{
		{SatID: "pipit-1", Station: "PEL", MaxElevDeg: 42.5},
		{SatID: "thistle", Station: "LAN", MaxElevDeg: 8.1},
		{SatID: "furrow-a", Station: "SQT", MaxElevDeg: 71},
	}
	want := []Pass{
		{SatID: "pipit-1", Station: "PEL", MaxElevDeg: 42.5},
		{SatID: "furrow-a", Station: "SQT", MaxElevDeg: 71},
	}
	got := Filter(passes, func(p Pass) bool { return p.MaxElevDeg > 10 })
	if !slices.Equal(got, want) {
		t.Errorf("Filter(passes, above 10 degrees) = %v, want %v", got, want)
	}
}

func TestFilterInts(t *testing.T) {
	got := Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
	if want := []int{2, 4, 6}; !slices.Equal(got, want) {
		t.Errorf("Filter([1 2 3 4 5 6], even) = %v, want %v", got, want)
	}
}

// TestFilterLeavesItemsAlone checks that Filter builds its result in a new slice.
func TestFilterLeavesItemsAlone(t *testing.T) {
	codes := []string{"YST", "PEL", "SQT", "LAN", "CWR"}
	Filter(codes, func(code string) bool { return code != "PEL" })
	if want := []string{"YST", "PEL", "SQT", "LAN", "CWR"}; !slices.Equal(codes, want) {
		t.Errorf("after Filter(codes, not PEL), codes = %q, want %q", codes, want)
	}
}

func TestFilterKeepsNothing(t *testing.T) {
	got := Filter([]string{"YST", "PEL", "SQT"}, func(code string) bool { return code == "XYZ" })
	if len(got) != 0 {
		t.Errorf(`Filter([YST PEL SQT], code == "XYZ") = %q, want an empty slice`, got)
	}
}
