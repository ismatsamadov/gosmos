`Push` gets an `any`: `x.(PassView)` gives the pass back, and appending it to `*h` grows the heap. `heap.Push` then moves it up to its place.

In `Pop`, take `old := *h`, keep its last pass, set `*h` to everything before that pass, and return the pass you kept.

If the very first pop is wrong, `Pop` removes the first pass. By the time your `Pop` runs, `heap.Pop` has moved the earliest pass to the end and the next earliest to the front, so the first pass is the wrong one to take.
