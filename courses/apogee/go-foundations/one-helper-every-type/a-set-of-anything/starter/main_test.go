package main

import (
	"cmp"
	"maps"
	"slices"
	"testing"
)

func TestAddHas(t *testing.T) {
	s := Set[string]{}
	s.Add("PEL")
	if !s.Has("PEL") {
		t.Errorf(`after Add("PEL"), Has("PEL") = false, want true`)
	}
	if s.Has("YST") {
		t.Errorf(`Has("YST") = true, want false: YST was never added`)
	}
}

func TestIntersectStations(t *testing.T) {
	thistle := Set[string]{"PEL": {}, "LAN": {}, "SQT": {}}
	pipit := Set[string]{"LAN": {}, "CWR": {}, "PEL": {}}
	got := thistle.Intersect(pipit)
	if want := []string{"LAN", "PEL"}; !slices.Equal(sorted(got), want) {
		t.Errorf("thistle.Intersect(pipit) = %q, want %q", sorted(got), want)
	}
	if got, want := sorted(thistle), []string{"LAN", "PEL", "SQT"}; !slices.Equal(got, want) {
		t.Errorf("after thistle.Intersect(pipit), thistle = %q, want %q", got, want)
	}
	if got, want := sorted(pipit), []string{"CWR", "LAN", "PEL"}; !slices.Equal(got, want) {
		t.Errorf("after thistle.Intersect(pipit), pipit = %q, want %q", got, want)
	}
}

func TestIntersectFrameNumbers(t *testing.T) {
	heardAtLAN := Set[int]{4182: {}, 4183: {}, 4185: {}}
	heardAtPEL := Set[int]{4183: {}, 4184: {}, 4185: {}}
	got := heardAtLAN.Intersect(heardAtPEL)
	if want := []int{4183, 4185}; !slices.Equal(sorted(got), want) {
		t.Errorf("heardAtLAN.Intersect(heardAtPEL) = %v, want %v", sorted(got), want)
	}
}

// sorted returns a set's values in order, for comparisons and messages.
func sorted[T cmp.Ordered](s Set[T]) []T {
	return slices.Sorted(maps.Keys(s))
}
