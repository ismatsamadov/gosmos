---
slug: old-code-new-lines
title: Old Code, New Lines
type: code_tests
position: 6
difficulty: 6
topics: logging
created: 2026-11-07T07:51:36Z
modified: 2026-11-07T07:51:36Z
---
The booking service you built a year ago, bookingd, writes one line for every request it answers, such as `req-000007 POST /bookings 201 236B 4ms`, through a `*log.Logger`. Nkem wants those lines in the same JSON as uplinkd's, without touching the code that writes them.

`slog.NewLogLogger` returns a `*log.Logger` that hands every line it is given to a handler as one record, at the level you choose, with the line's text as the message. Old code calls `Printf` as before:

```go
pl := slog.NewLogLogger(h, slog.LevelWarn)
pl.Printf("%d passes moved", 3)
// {"level":"WARN","msg":"3 passes moved"}
```

## Your task

Make `accessLog` return a `*log.Logger` that hands every line to `h` as an INFO record, instead of writing it to standard error.

> [!NOTE]
> The values stay inside the message, so a search finds the line by its text, not by key: a bridge until bookingd logs with slog itself.
