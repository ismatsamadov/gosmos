Only the goroutine inside `read` sends on `out`, so only it knows when the last line has gone: the close belongs there, not after the `go` statement.

A deferred call runs as its function returns. Make `defer close(out)` the first line of the goroutine's function, before the scanner, so the channel closes however the loop ends.

If the test still says read did not close its channel, the `defer` is not inside `go func() { ... }()`. If it panics with `send on closed channel`, the close is in `read` itself, which returns before the first send.
