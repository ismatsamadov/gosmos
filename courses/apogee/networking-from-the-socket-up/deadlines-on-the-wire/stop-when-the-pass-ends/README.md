---
slug: stop-when-the-pass-ends
title: Stop When the Pass Ends
type: code_tests
position: 6
difficulty: 8
topics: context, networking, time
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Ingrid's rule: at LOS, when the satellite sets out of view, dishd stops waiting on Dish 1 at once, not 5 seconds later.

A pass runs under a `context.Context` that ends at LOS. `context.AfterFunc(ctx, f)` runs `f` in a goroutine when `ctx` ends; its `stop` cancels `f` if called first, or reports `false` without waiting for `f`. From `f`, a read deadline of now wakes a waiting read:

```go
stop := context.AfterFunc(ctx, func() {
	fmt.Println("LOS: the pass is over")
})
defer stop() // done before LOS: nothing prints
```

## Your task

1. In `recvCtx`, if `ctx` has already ended, return its error.
2. Register `context.AfterFunc` to set `c.conn`'s read deadline to now, then close a channel. Call `c.recv()`, then `stop`; if it reports `false`, wait for the channel.
3. If the read failed and `ctx` has ended, return `ctx.Err()`, else what `recv` returned.

> [!WARNING]
> Skip `stop` or that wait, and a late deadline cuts the next read.
