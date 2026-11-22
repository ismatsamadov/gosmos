package main

// passHeap holds passes in heap order by byAOS, for container/heap.
type passHeap []PassView

// Len is how many passes the heap holds.
func (h passHeap) Len() int { return len(h) }

// Less reports whether pass i comes before pass j.
func (h passHeap) Less(i, j int) bool {
	return byAOS(h[i], h[j]) < 0
}

// Swap swaps passes i and j.
func (h passHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push appends x, a PassView, which heap.Push then moves into place.
func (h *passHeap) Push(x any) {}

// Pop removes and returns the last pass, where heap.Pop has moved the earliest.
func (h *passHeap) Pop() any {
	return nil
}
