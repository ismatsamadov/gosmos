`l.Enabled(ctx, slog.LevelDebug)` returns true only when `l`'s handler would write a DEBUG record right now, and it costs a comparison, not a dump.

Start `logFrame` with an `if` on that call: when it is false, `return` before anything else runs, and keep the `DebugContext` line as it is for the true case.

If the test counts calls to `dump` at INFO, the arguments are still worked out before the handler drops the record: the check must come before the `DebugContext` call, not inside it.
