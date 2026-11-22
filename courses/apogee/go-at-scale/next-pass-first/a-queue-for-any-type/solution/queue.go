package main

import "container/heap"

// Queue is a priority queue: Pop returns the smallest value by cmp.
type Queue[T any] struct{ h *heapOf[T] } // heapOf implements heap.Interface for a slice and cmp

// heapOf is the slice behind a Queue, kept in heap order by cmp for container/heap.
type heapOf[T any] struct {
	items []T
	cmp   func(a, b T) int
}

// Len is how many values the heap holds.
func (h *heapOf[T]) Len() int { return len(h.items) }

// Less reports whether value i comes before value j.
func (h *heapOf[T]) Less(i, j int) bool {
	return h.cmp(h.items[i], h.items[j]) < 0
}

// Swap swaps values i and j.
func (h *heapOf[T]) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

// Push appends x, which heap.Push then moves into place.
func (h *heapOf[T]) Push(x any) {
	h.items = append(h.items, x.(T))
}

// Pop removes and returns the last value, where heap.Pop has moved the smallest.
func (h *heapOf[T]) Pop() any {
	last := h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	return last
}

// NewQueue returns an empty queue that pops the smallest value by cmp first.
func NewQueue[T any](cmp func(a, b T) int) *Queue[T] {
	return &Queue[T]{h: &heapOf[T]{cmp: cmp}}
}

// Push adds v to the queue.
func (q *Queue[T]) Push(v T) { heap.Push(q.h, v) }

// Pop removes and returns the smallest value, or the zero value and false when the queue is empty.
func (q *Queue[T]) Pop() (T, bool) {
	if q.h.Len() == 0 {
		var zero T
		return zero, false
	}
	return heap.Pop(q.h).(T), true
}

// Len is how many values wait in the queue.
func (q *Queue[T]) Len() int { return q.h.Len() }
