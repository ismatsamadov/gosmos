For `newPassLog`, make the JSON handler on `w` with no options, give it to `stamp` together with `now`, and hand what `stamp` returns to `slog.New`.

For `passLine`, make the three attributes both lines share, with `slog.String`, `slog.Float64` and `msAttr`; then choose the level and the message by whether `err` is nil, adding `slog.Any("err", err)` only for a failure.

If the test reports the Mac's own clock as the time, the JSON handler isn't inside `stamp`. If `wait_ms` is missing, or a test shows the wait in nanoseconds, the wait went in without `msAttr`.
