Start `srv.Serve(ln)` in a goroutine that sends its result on a channel with room for one, so it never waits for anyone to receive it. After that, call `b.Refresh(ctx)` once; when it fails, the Board keeps the view it had, so carry on.

The rest is a `for` loop around a `select` of three cases: a tick from `ticks` refreshes again; a result from the serving goroutine means `Serve` failed, so return it; `ctx.Done()` means it is time to stop.

To stop, make a context with `context.WithTimeout` from `context.Background()`, not from `ctx`, which is already done, and pass it to `srv.Shutdown`. Then receive `Serve`'s result: `http.ErrServerClosed` is expected, so return nil.
