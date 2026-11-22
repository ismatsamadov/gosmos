Keep a `body` slice outside a `for` loop. Each turn reads one size line and, unless the size is 0, one chunk: `make` a slice of the size plus 2 and fill it with `io.ReadFull(br, ...)`.

The last two bytes of that slice must be `"\r\n"`; the rest is the chunk's data, which you `append` to the body. `errors.Is(err, errLineTooLong)` tells a size line that doesn't fit apart from other errors, which you return as they are.

If a test gets an empty body back, the loop never runs or the data isn't appended. If the 4,096-byte chunk fails, the data is read with `readLine` or a single `Read` rather than `io.ReadFull`.
