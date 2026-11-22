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
}

// Has reports whether v is in the set.
func (s Set[T]) Has(v T) bool {
	return false
}

// Intersect returns a new set of the values in both s and other.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	return Set[T]{}
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
