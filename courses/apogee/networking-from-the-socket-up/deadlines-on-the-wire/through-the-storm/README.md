---
slug: through-the-storm
title: Through the Storm
type: code_tests
position: 7
difficulty: 8
topics: context, errors, networking, time
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Friday 14 April. Sunniva runs Tuesday's storm again on the simulator in the hut: a slow line, a silent one, one that stops reading, and a pass that ends while dishd waits. "Storm, take two. Coffee ready."

`Do` is the one call per frame: it sends with the write deadline, then waits with `recvCtx`, so a slow line still answers, a dead one fails within 5 seconds, and a pass's end ends the wait. After a failed send, there is no answer to wait for.

```go
_, err := c.Do(ctx, Frame{Type: TypePing})
if err != nil {
	log.Println(describe(c.Name, c.Timeout, err))
	// YST-1: no answer in 5s
}
```

## Your task

`Do` returns zero values. Send `f` with `c.send`, returning its error at once if it fails, then return what `recvCtx(ctx, c)` returns.

> [!TIP]
> Reading after a failed send only adds 5 more seconds to the failure, and hides the error that says what broke.
