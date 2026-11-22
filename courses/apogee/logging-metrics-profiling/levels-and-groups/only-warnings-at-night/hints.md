The options are made on the first line of `main`, and `slog.HandlerOptions` has a field for the lowest level its handler writes: `Level`.

Give `Level` the constant for WARN, `slog.LevelWarn`, beside `ReplaceAttr`, and remove the TODO comment. The five calls below it stay as they are.

If Submit still shows `level=INFO msg="pass started"`, `Level` is unset or set to INFO; the DEBUG lines never showed, because INFO is the default minimum.
