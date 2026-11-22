---
slug: contexts-all-the-way-down
title: Contexts All the Way Down
type: code_tests
position: 5
difficulty: 10
topics: context, http
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
In Wednesday's rehearsal, one `POST /bookings` stuck behind a slow disk held the old process for the whole twenty seconds, and it was still waiting when the shutdown gave up.

Every request's context is made from the server's base context, which `srv.BaseContext` can supply. If `run` owns that base, it can cancel it once `Shutdown` has run out of time: requests that would finish still do, and the stuck ones stop waiting, because a context made from another ends with it:

```go
base, cancel := context.WithCancel(
	context.Background())
child := context.WithValue(base, idKey{}, 42)
cancel()
fmt.Println(child.Err()) // context canceled
```

## Your task

In `run`, set `srv.BaseContext` to return a context that `run` cancels only when `stop` fails. Then cancel it, `Close` the server, and return the error wrapped as `shutdown: <error>`.

> [!WARNING]
> Cancel the base only after `Shutdown` gives up: canceled sooner, it cuts short requests that would have finished.
