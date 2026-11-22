---
slug: draining-not-down
title: Draining, Not Down
type: code_tests
position: 4
difficulty: 9
topics: http, time
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Yusra's customer guide tells customers' programs to check `/health` before they send a batch of bookings. During a restart they should wait, not start a batch that the restart cuts in half.

So the service drains before it stops: for `t.Drain`, `/health` answers `503` with `{"status":"draining"}` while every other route still works, and only then does the shutdown begin. The flag that `run` sets and every health check reads is an `atomic.Bool`, safe to share between goroutines:

```go
var paused atomic.Bool
paused.Store(true)
fmt.Println(paused.Load()) // true
```

## Your task

1. In `run`, once `ctx` ends, set `app.Service.draining`, keep serving for `t.Drain` with `time.Sleep`, then stop.
2. Make `health` answer 503 with `{"status":"draining"}` while the flag is set.
