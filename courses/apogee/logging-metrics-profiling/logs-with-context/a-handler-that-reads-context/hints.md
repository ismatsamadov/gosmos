`passFrom(ctx)` and `requestIDFrom(ctx)` each tell you whether the context holds an ID: the first with its `bool`, the second by returning `""` when there is none.

For each one the context holds, call `r.AddAttrs` with `slog.String` and its key, `pass_id` or `request_id`, before the `return` line passes `r` to `h.Handler.Handle`.

If the test finds no `pass_id`, the attributes were added after `r` was passed on, or to a copy that was never passed; `AddAttrs` changes `r` itself, so add to `r` first and pass that same `r`.
