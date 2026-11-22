---
slug: dont-build-unread-lines
title: Don't Build Unread Lines
type: code_tests
position: 6
difficulty: 7
topics: context, functions, logging
created: 2026-11-10T15:36:37Z
modified: 2026-11-10T15:36:37Z
---
Thursday, on the lab bench: each frame the radio receives, a block of bytes from the satellite, goes into a DEBUG record as hex. Making that hex costs more than flying the pass, and at INFO nobody reads it.

Go works out a call's arguments before making the call, so `l.Debug(..., dump(frame))` builds the hex even when the handler then drops the record. `Logger.Enabled` asks the handler first whether a record at a level would be written; when it says no, skip the work:

```go
if l.Enabled(ctx, slog.LevelDebug) {
	l.DebugContext(ctx, "lines sent",
		"text", strings.Join(lines, " | "))
}
```

## Your task

Make `logFrame` call `dump` only when `l` would write a DEBUG record, and write nothing at all otherwise.

> [!NOTE]
> `Enabled` asks the handler each time, so on a logger with a `LevelVar` the answer changes the moment the night shift turns the station up.
