Only the handler changes: the line in `newLogger` that calls `slog.NewTextHandler` is the one to edit.

`slog.NewJSONHandler` takes the same two arguments as `slog.NewTextHandler`, the writer and the options. Keep both as they are, so `dropTime` still leaves the time out.

If the test quotes a line that starts with `level=`, the text handler still writes it; a line from the JSON handler starts with `{`.
