`l.LogAttrs` takes a context first, and `context.Background()` is the one to give it; then the level, `slog.LevelInfo`, the message, and the attributes, separated by commas.

Make each attribute with the function for its type, its key first: `slog.String` for `p.ID`, `slog.Float64` for `p.MaxElevDeg` and `slog.Int` for `lines`. Remove the `_ = context.Background` line once you call it.

If the test says `"max_elev_deg" = "16.6"`, the elevation went in as text, through `fmt` or `slog.String`: give the float itself to `slog.Float64`.
