Give `slog.HandlerOptions` a `ReplaceAttr` that returns `slog.Time(slog.TimeKey, w.Now())` for the top-level `slog.TimeKey` attribute, the one that comes with no groups, and add `program` and `site` to each station's logger once, with its `With` method. Hand times and numbers over as `time.Time`, `int` and `float64` values: the handler writes them.

Write each log through a `bufio.Writer` into a temporary file in `logs/`, behind a mutex so a record is one `Write`. Once the week's passes and deliveries are done, stopped or not, flush it, close it and rename it into place, before `rehearse` returns.

Log "pass started" right after the AOS line, the end record right after the LOS line unless the run stopped, and "line lost" where a welcomed call's read fails. Keep each failed call's error on its pass (the dial, the banner, the `HELLO`, a read, a NAK, `BYE stowed`) for `err`, count `frames` before a TLM line is parsed, and count `commands` where an answer settles a command this pass sent.
