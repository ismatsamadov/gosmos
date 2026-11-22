---
slug: stop-every-stage
title: Stop Every Stage
type: code_tests
position: 3
difficulty: 8
topics: channels, context, select
created: 2026-11-10T15:39:51Z
modified: 2026-11-11T12:51:42Z
---
In the drill, the Board read three frames and left, and the parse stage's goroutine sat blocked on its next send for the rest of the night, holding its feed open.

A goroutine blocked on a send nobody receives, or on a receive nobody sends to, waits forever. So in a stage that can stop, every send and every receive is a `select` beside `<-ctx.Done()`, and the goroutine returns when the context ends, closing its channel on the way out.

```go
select {
case alarms <- kmh:
case <-ctx.Done():
	return // nobody is reading any more
}
```

## Your task

`parse` ranges over `in` and sends with a plain `out <- f`. Make the receive and the send each a `select` that also waits on `ctx.Done()`, and return when the context is done.

> [!TIP]
> In a `select`, `case h, ok := <-in:` receives a line, and `ok` is false once `in` is closed.
