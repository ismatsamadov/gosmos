The waiting checks stop only when their context is canceled. Which goroutine knows when to cancel it, and with what?

Make the child context with `context.WithCancelCause` at the top of `runAll`, and defer `cancel(nil)` to release it. Its cancel function takes the failing check's error, and `context.Cause` reads it back once `wg.Wait()` returns.
