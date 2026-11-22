---
slug: wait-for-room
title: Wait for Room
type: code_tests
position: 3
difficulty: 8
topics: channels, context, select
created: 2026-11-11T12:51:42Z
modified: 2026-11-12T08:43:35Z
---
Wednesday 14 March. Varsk Island (VSK), a new station at a North Pacific fishing harbor, takes the passes passplan 2 gives its dish, its contacts, over a thin line. The first draft kept every contact it couldn't send yet in a slice: 61,000 by this morning.

Its dish loader, the program that loads each contact for the dish, reads them from a queue with room for 256. When the queue is full, the sender waits for room, or for the run to be canceled: that is back-pressure, and it keeps memory bounded however slow the line.

```go
for _, f := range frames {
	select {
	case toBoard <- f: // waits while it's full
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

## Your task

`feed` hands the contacts to a goroutine and returns at once. Instead, in a `select` beside `<-ctx.Done()`, send each one into `q` in order, and return `nil` once all are in, or `ctx.Err()`.

> [!WARNING]
> Returning early looks fast, but the contacts still wait in memory, in a goroutine nobody can stop.
