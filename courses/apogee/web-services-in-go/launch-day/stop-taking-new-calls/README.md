---
slug: stop-taking-new-calls
title: Stop Taking New Calls
type: code_tests
position: 0
difficulty: 9
topics: context, http, time
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Monday 24 January. One line of Ingrid's launch checklist has never been tried: "Restart bookingd during a pass." The booking service's program, bookingd, will restart while customers' programs are in the middle of requests.

`srv.Close()` cuts every connection at once, requests in flight included. `srv.Shutdown(ctx)` stops gracefully: it closes the listener, so no new connection gets in, then waits for the requests in flight to finish or for `ctx` to end, whichever comes first, and returns `ctx.Err()` if `ctx` won. A context with a timeout ends on its own:

```go
ctx, cancel := context.WithTimeout(
	context.Background(), 2*time.Second)
defer cancel()
<-ctx.Done()
fmt.Println(ctx.Err())
// context deadline exceeded
```

## Your task

`stop` calls `srv.Close()`. Call `srv.Shutdown` instead, with a context that ends after `within`. A `within` of 0 means no limit: wait for every request in flight, however long it takes.
