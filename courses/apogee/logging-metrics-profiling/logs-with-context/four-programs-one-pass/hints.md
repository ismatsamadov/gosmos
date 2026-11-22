Build it from the inside out: `slog.NewJSONHandler(w, nil)` first, then `stamp` around it with `now`, then a `passHandler` around what `stamp` returns, as a struct literal with that handler inside.

Make the logger with `slog.New` on the `passHandler`, then call `With` on it with the two keys and their values, `program` first, and return the logger `With` returns.

If the test reports the Mac's own clock, the JSON handler isn't inside `stamp`; if `pass_id` is missing, the records skip `passHandler`; if `program` is missing, `With` was called on a logger that isn't the one returned.
