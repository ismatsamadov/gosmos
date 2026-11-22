---
slug: first-failure-stops-all
title: First Failure Stops All
type: code_tests
position: 5
difficulty: 9
topics: context, errors, goroutines
created: 2026-11-11T12:51:42Z
modified: 2026-11-11T12:51:42Z
---
Thursday 15 March. When a satellite's elements, the numbers that place it on its orbit, are missing, its prediction fails and the other seventy-nine run on for nothing. Nkem wants passplan 2 to stop at its first failure, as the pre-pass checks have since 2027.

A group is slots and a `sync.WaitGroup` together. Its first failure cancels the jobs' context, with the error as the cause, and stays for `Wait`. `sync.Once` keeps it exactly once: `once.Do(f)` runs `f` the first time only. Next week shows what else it promises.

```go
var once sync.Once
for _, c := range []string{"TIR", "VSK"} {
	once.Do(func() { fmt.Println("first:", c) })
}
// first: TIR
```

## Your task

1. Write `fail(err)`: in `g.once.Do`, keep `err` in `g.err` and pass it to `g.cancel`. In `Go`'s goroutine, hand a job's error to `fail`.
2. Once `Go` has a slot, if `context.Cause(g.ctx)` isn't `nil`, give the slot back, `fail` with that cause, and run nothing.

> [!TIP]
> A job `Go` refuses is a failure too: its cause is the group's first failure, or why its caller canceled.
