The `n` that the inner `Write` returns is the number to add: convert it with `int64(n)` and add it to `c.n`.

Add it before you return, on every path: a write that fails partway has still put `n` bytes in the file. Return `n` and `err` exactly as they came, without wrapping the error.

If the test says the count is 1000 after a short write, you counted only writes that succeeded; if it says 1085, you added `len(p)` instead of `n`.
