---
slug: nothing-left-behind
title: Nothing Left Behind
type: code_tests
position: 7
difficulty: 10
topics: context, defer, errors, goroutines, select
created: 2026-11-10T15:36:39Z
modified: 2026-11-11T17:33:39Z
---
Friday 9 February, 10:00: the fix goes to every station. At Thursday's review Nkem had one question: "What does the worker see if this pass never ends?"

A `select` waits on several channel operations and takes whichever can go first: a free worker, or the end of the pass. With no goroutine, nothing outlives `runPass`, and since Thursday the radio's `fly` takes the pass's context and stops at LOS with `ErrNoAck`:

```go
select {
case queue <- line:
	fmt.Println("queued")
case <-ctx.Done():
	fmt.Println("too late:", ctx.Err())
}
```

## Your task

Rewrite `runPass` with no goroutine:
1. Take a worker in a `select` on `u.workers` and `ctx.Done()`; if the pass ends first, return `ErrNoWorker` wrapped as `pass <ID>: ...`.
2. Give the worker back with `defer`.
3. Call `u.fly` with `ctx` and the pass's buffer, and return its error wrapped the same way, or nil.

> [!NOTE]
> The test replays the night of 6–7 February and checks, with your `inUse` and `leftBehind` from `memory.go`, that goroutines and the heap come back to where they started.
