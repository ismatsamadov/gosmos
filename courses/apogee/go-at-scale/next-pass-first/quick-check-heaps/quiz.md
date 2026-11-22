## Who calls `passHeap`'s `Pop` method in this code?

- [x] `heap.Pop`, after it has moved the earliest pass to the end of the slice
  > The package does the heap's work and calls your method only to take the last element off the slice.
- [ ] Your code, each time it needs the next pass
  > Called directly, your method takes the last element, which is the earliest only after `heap.Pop` has moved it there.
- [ ] `heap.Push`, to make room for the new pass
  > `heap.Push` calls your `Push`, then moves the new pass up to its place; it never calls `Pop`.

## Why does your `Pop` take the last element?

- [ ] A heap keeps its slice sorted from the latest pass to the earliest
  > A heap isn't sorted: only index 0 is sure to hold the smallest value, and the rest are in heap order.
- [ ] Taking the first element would copy the whole slice
  > Reslicing from 1 copies nothing; the trouble is that the first element is no longer the one `heap.Pop` wants removed.
- [x] `heap.Pop` has already swapped the smallest value there
  > It swaps the first and last elements, restores the order of the rest, then calls your `Pop` to take the last one off.

## How does the work of one `heap.Pop` grow with the queue?

- [ ] With n: a pop from 20,000 passes looks at every one of them
  > That's a search through every pass; a pop never looks at most of the heap's values.
- [x] With log n: a heap of 20,000 passes is about fifteen levels deep
  > A pop moves one value down from the top, a level at a time, so doubling the queue adds only one level.
- [ ] Not at all: the smallest value is always at index 0
  > Reading index 0 costs nothing, but removing it means moving another value down to its place.
