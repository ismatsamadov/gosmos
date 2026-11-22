`net.Pipe()` returns both ends at once, and `pipePair` returns them in the same order.

Register one function with `t.Cleanup` that closes both ends. It runs when the test ends, after every write and read the test makes.

If the test meets `io: read/write on closed pipe` while it still uses the pair, the ends were closed too early, by a `defer`. If a `Read` after the test gives anything but `io.ErrClosedPipe`, one end or both were never closed.
