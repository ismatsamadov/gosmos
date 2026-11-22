Pass `wg.Go` a function that calls `handle(c)`: the handler runs in a goroutine of its own, counted by `wg`, and the loop is back at `Accept` at once.

When `Accept` fails, check `errors.Is(err, net.ErrClosed)` first and return `nil` for it; return any other error unchanged.

If the watch port accepts no second watcher, the handler still runs in the loop's own goroutine. If `wg.Wait()` returns while handlers are still running, they were started with `go` instead of `wg.Go`.
