---
slug: every-ten-minutes
title: Every Ten Minutes
type: code_tests
position: 3
difficulty: 8
topics: context, select, time
created: 2026-11-11T16:51:57Z
modified: 2026-11-12T08:43:36Z
---
Wednesday 25 April. passplan 2 replans every ten minutes, but its draft waits on a new `time.After` each round, counted from the end of the last replan, so every slow replan pushes the rest of the day later.

A `time.Ticker` sends the time on its channel `C` every `d`, on a beat fixed when it is made: one of ten minutes ticks at 10, 20 and 30 minutes, however long the work between ticks takes, as long as it takes less than `d`. `Stop` ends the ticks:

```go
tick := time.NewTicker(30 * time.Second)
defer tick.Stop()
for range 3 {
	at := <-tick.C // 30, 60 and 90 s in
	fmt.Println("Board refresh", at.UTC())
}
```

## Your task

`every` calls `f` once and returns. Make a ticker of `d` and `defer` its `Stop`, then call `f` with each tick's time until the context ends, in a `select` on the ticks and `ctx.Done()`. Return as soon as the context ends.

> [!WARNING]
> Inside a `select`, `break` leaves only the `select`, and the loop around it goes on: `return` from `every` instead.
