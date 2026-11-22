---
slug: a-line-in-the-log
title: A Line in the Log
type: code_tests
position: 4
difficulty: 7
topics: http, time
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Yusra pastes a request ID from Cirrine's email, `req-000007`, and wants the one line that says what happened to it. That line comes from an access log: one line per request, written when its handler has returned and the status and size are known.

`log.New(w, "", 0)` makes a logger with no prefix and no date, safe for many requests at once. The service's clock is a function, `now`, so a test can move time by hand:

```go
l := log.New(os.Stdout, "", 0)
start := now()
checkWind()
took := now().Sub(start)
l.Print(took.Round(time.Millisecond))
// 1.204s
```

## Your task

`logRequests` calls `next` and logs nothing. Wrap `w` in a `statusRecorder`, call `next` with it, then log one line: the request ID, the method, the path without its query, the status (200 if the handler wrote nothing, as the server sends), the bytes, and the time taken, rounded to the millisecond:

```text
req-000007 POST /bookings 201 236B 4ms
```

> [!WARNING]
> Log after `next` returns: before that, the status is still 0.
