package main

import "fmt"

// Ring keeps the last items pushed into it, dropping the oldest.
type Ring[T any] struct {
	items []T // the items, never more than size of them
	size  int // how many items the ring keeps
	next  int // where Push writes once the ring is full
}

// NewRing returns an empty ring that keeps the last n items; n must be at least 1.
func NewRing[T any](n int) *Ring[T] {
	return &Ring[T]{items: make([]T, 0, n), size: n}
}

// Push adds v, dropping the oldest item if the ring is full.
func (r *Ring[T]) Push(v T) {
	if len(r.items) < r.size {
		r.items = append(r.items, v)
		return
	}
	r.items[r.next] = v
	r.next = (r.next + 1) % r.size
}

// Items returns the items in a new slice, oldest first.
func (r *Ring[T]) Items() []T {
	out := make([]T, 0, len(r.items))
	out = append(out, r.items[r.next:]...)
	return append(out, r.items[:r.next]...)
}

func main() {
	volts := NewRing[float64](3)
	for _, v := range []float64{7.91, 7.88, 7.86, 7.84} {
		volts.Push(v)
	}
	fmt.Println(volts.Items())
}
