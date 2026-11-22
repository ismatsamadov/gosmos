---
slug: milliseconds-not-nanoseconds
title: Milliseconds, Not Nanoseconds
type: code_tests
position: 5
difficulty: 6
topics: logging, time
created: 2026-11-07T07:51:36Z
modified: 2026-11-07T07:51:36Z
---
On the lab bench at Port Ellery, where every change runs before it reaches a station, Ingrid reads the new failure line's wait, 570000000000, and takes it for a phone number. It is the 01:55 pass's 9 minutes 30 seconds, from its preparation to its end, counted in nanoseconds.

The JSON handler writes a `time.Duration` as its count of nanoseconds. Apogee logs a duration as whole milliseconds instead, under a key that names the unit. `Duration.Milliseconds` returns them as an `int64`, dropping any fraction:

```go
d := 4*time.Second + 250*time.Millisecond
fmt.Println(int64(d))         // 4250000000
fmt.Println(d.Milliseconds()) // 4250
```

## Your task

Make `msAttr` return an `slog.Int64` attribute: its key is `key` with `_ms` added, its value `d.Milliseconds()`. Then `msAttr("wait", d)` gives `wait_ms`, and the 01:55 pass's wait reads `570000`.

> [!NOTE]
> A duration under a millisecond rounds down to 0: 999µs gives `0`.
