---
slug: a-line-that-says-why
title: A Line That Says Why
type: code
position: 0
difficulty: 6
topics: logging
created: 2026-11-07T07:51:36Z
modified: 2026-11-07T07:51:36Z
---
Wednesday 8 November, 02:48, at home. The pager woke you two minutes ago: Lantern Atoll, the island station near the equator, is failing one pass in five, and the log of its uplinkd, the program that hands each pass's commands to the dish, holds three lines that each end in `error`.

A structured log is made of records instead of free text. A `*slog.Logger`, from `log/slog`, writes each record as a message that never changes, then key-value pairs, its attributes, that carry what does. Each method is named for the record's log level, how much it matters, such as `Info` or `Error`:

```go
logger.Info("lines queued", "lines", 1)
// level=INFO msg="lines queued" lines=1
```

## Your task

Replace `log.Print("error")` in `main` with one `logger.Error` call: the message `pass failed` and the attribute `site` set to `LAN`. The program must print:

```text
level=ERROR msg="pass failed" site=LAN
```

> [!NOTE]
> `log.Print` writes the old line, `2028/11/08 02:46:24 error`, to standard error, which Submit doesn't read. `dropTime`, in `clock.go`, leaves the time out of the new line so it fits your screen.
