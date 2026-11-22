`slog.Int64` takes a key and an `int64`, and `d.Milliseconds()` is already an `int64`: no conversion is needed.

Build the key by adding `"_ms"` to `key` with `+`, so that `msAttr("ack", d)` gives `ack_ms`.

If the test shows `"wait":570000000000`, the attribute is still a duration; if it shows `"wait_ms":"570000"`, the number went in as text.
