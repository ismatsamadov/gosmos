Make the child with `context.WithTimeout(ctx, limit)`, `defer` its `cancel` on the next line, and hand the child, not the caller's `ctx`, to `plan`.

`errors.Is(err, context.DeadlineExceeded)` picks out a deadline: then return `nil` and a new error from `fmt.Errorf` that prints `limit` with `%v` and ends in `: %w`. Otherwise keep what `plan` returns. Add `errors` and `fmt` to the imports.

If the 31 s run still returns at 31s, `plan` isn't getting the child. If the run canceled at 10 s gets the timeout's message, wrap the deadline only, not every error. If the run's context is still live, `cancel` isn't called when `planWithin` returns.
