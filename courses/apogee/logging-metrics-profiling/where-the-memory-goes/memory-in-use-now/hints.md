Two lines change: add a call to `runtime.GC()` at the top of `inUse`, before `ReadMemStats`, and return a different field of `m`.

`HeapAlloc` is the bytes of the heap's blocks that haven't been freed. Right after a collection, the only blocks left are the ones something can still reach: the memory in use.

If the test says the number never came back down, either `TotalAlloc` is still returned, or `HeapAlloc` was read without a collection first, so the let-go buffers still counted as garbage.
