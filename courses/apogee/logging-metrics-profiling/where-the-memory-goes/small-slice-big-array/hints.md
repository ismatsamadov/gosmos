The starter's slice expression already picks the right bytes, and its `max` keeps it from going below 0 for a short buffer. What it returns is a window onto `frames`; what you need is the same bytes in an array of their own.

`bytes.Clone(b)` returns a new slice holding a copy of `b`'s bytes. Remove the `_ = bytes.Clone` line, and wrap the slice expression in `bytes.Clone`.

If the test says the tails keep megabytes in use, the result still shares the buffer's array, however few bytes it shows: only a copy lets the buffer go.
