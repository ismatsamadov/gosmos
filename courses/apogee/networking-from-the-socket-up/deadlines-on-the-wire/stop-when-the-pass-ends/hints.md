`ctx.Err()` is `nil` while the context is on; once it has ended, it says why: `context.Canceled` or `context.DeadlineExceeded`.

The function you give `AfterFunc` sets the read deadline of `c.conn` to `time.Now()`, an instant already reached, then closes a channel made with `make(chan struct{})`. Call the `stop` it returns as soon as `c.recv()` has returned; when `stop` reports `false`, receive from the channel before going on.

If a later read fails after an earlier pass's context ended, `stop` was never called. If `recvCtx` returns while the deadline set at LOS is still on its way, wait for the channel whenever `stop` reports `false`. If the first case gives `i/o timeout` after 20ms, return `ctx.Err()`, not the read's own error, when the context has ended.
