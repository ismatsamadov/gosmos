Make the base with `context.WithCancel(context.Background())`, never from `ctx`: `ctx` ends at the signal, long before the shutdown gives up. `BaseContext` is a function that takes the listener and returns the base.

Keep the `cancel` function, with a `defer` for the paths that never need it, and call it only when `stop` returns an error.

After canceling, `srv.Close()` cuts the connections left, and wrapping the error with `%w` after `shutdown: ` keeps `context.DeadlineExceeded` findable with `errors.Is`. If the request that would have finished gets 503, the base was canceled, or taken from `ctx`, too early.
