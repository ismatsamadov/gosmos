---
slug: stop-without-cutting
title: Stop Without Cutting
type: code_tests
position: 7
difficulty: 10
topics: context, goroutines, networking
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T19:06:41Z
---
Monday night. Nkem will install the station server's first update during Tuesday's passes, and stopping it must never cut an upload in the middle.

A clean stop goes in order: no new connections; those waiting between requests closed; every feed ended; each request in flight left to finish its answer, after which `send` hangs up. `s.conns` marks each open connection true while it is busy, and the given `closeConns` closes the waiting ones, or all. A `select` takes whichever comes first:

```go
select {
case <-passOver:
	fmt.Println("LOS")
case <-ctx.Done():
	fmt.Println(ctx.Err())
}
```

## Your task

`Shutdown` only closes the listener. Then close the waiting connections, end the feeds with `s.Store.Close()` and wait for `s.wg`, returning nil; if `ctx` ends first, close everything, wait, and return `ctx.Err()`.

> [!WARNING]
> Close the waiting connections before you wait: each one waits for its next request, and `Wait` would wait with it.
