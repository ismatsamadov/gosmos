package main

import (
	"fmt"
	"maps"
	"slices"
)

// Set holds each value at most once; its struct{} values take no memory.
type Set[T comparable] map[T]struct{}

// Add puts v in the set.
func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

// Has reports whether v is in the set.
func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

// Intersect returns a new set of the values in both s and other.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	both := Set[T]{}
	for v := range s {
		if other.Has(v) {
			both.Add(v)
		}
	}
	return both
}

func main() {
	thistle, pipit := Set[string]{}, Set[string]{}
	for _, code := range []string{"PEL", "LAN", "SQT"} {
		thistle.Add(code)
	}
	for _, code := range []string{"LAN", "CWR", "PEL"} {
		pipit.Add(code)
	}
	both := thistle.Intersect(pipit)
	fmt.Println("heard both:", slices.Sorted(maps.Keys(both)))
}
