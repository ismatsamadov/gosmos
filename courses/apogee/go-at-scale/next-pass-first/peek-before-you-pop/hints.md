The smallest value is the first element of `q.h.items`, at index 0: `heap.Init`, `heap.Push` and `heap.Pop` all keep it there.

Check `q.h.Len()` before you index the slice. For an empty queue, keep the stub's `var zero T` and false; otherwise return `q.h.items[0]` and true.

If `Len` drops after a `Peek`, the method pops instead of reading. If an empty queue panics, the length check is missing.
