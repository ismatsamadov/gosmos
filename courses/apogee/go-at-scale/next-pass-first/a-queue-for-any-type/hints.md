`NewQueue` returns `&Queue[T]{h: ...}`, with a `heapOf[T]` whose `cmp` field is the `cmp` it was given. `Len` returns `q.h.Len()`.

`Push` is one call to `heap.Push` with `q.h` and `v`. `Pop` checks `q.h.Len()` first: an empty queue keeps the stub's answer, and any other returns `heap.Pop(q.h).(T)` and true.

If `Pop` on an empty queue panics, it calls `heap.Pop` without checking the length: on an empty heap, `heap.Pop` reaches past the end of the slice.
