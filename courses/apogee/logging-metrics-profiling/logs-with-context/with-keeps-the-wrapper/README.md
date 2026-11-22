---
slug: with-keeps-the-wrapper
title: With Keeps the Wrapper
type: code_tests
position: 2
difficulty: 8
topics: interfaces, logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
Wednesday. After `l.With("program", "uplinkd")` on the lab bench, the pass ID vanishes from every one of Lantern Atoll's lines.

`Logger.With` calls its handler's `WithAttrs` and builds the new logger on the handler that comes back; `WithGroup` works the same way. `passHandler`'s two methods hand back the inner handler's result, a handler with no `passHandler` around it, so every logger made with `With` skips yours. A wrapper returns a new wrapper around the result instead, as `stamp` did in Logs a Machine Can Read:

```go
func (h quiet) WithAttrs(
	as []slog.Attr) slog.Handler {
	return quiet{h.Handler.WithAttrs(as)}
}
```

## Your task

Make `passHandler`'s `WithAttrs` and `WithGroup` return a `passHandler` around what the inner handler returns, so loggers made with `With` and `WithGroup` still add the context's IDs. The test also runs `testing/slogtest`, which checks that a handler keeps the rules slog's own handlers keep.

> [!NOTE]
> A group opened with `WithGroup` never closes: every attribute after it lands inside, the IDs your `Handle` adds included. That is why Apogee writes its groups with `slog.Group`.
