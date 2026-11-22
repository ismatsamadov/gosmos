Swap `context.WithCancel` for `context.WithCancelCause`. Its cancel takes an error, so the deferred call becomes `defer cancel(nil)`.

Pass that cancel itself to every `readFeed`, in place of the small function that throws the error away.

Once `newest` returns, every stage has stopped: if `context.Cause(ctx)` isn't `nil`, return `nil` and that cause. If the test reports `context canceled`, the code returns `ctx.Err()` instead.
