---
slug: stop-twice-close-once
title: Stop Twice, Close Once
type: code_tests
position: 2
difficulty: 8
topics: channels, functions, goroutines
created: 2026-11-11T12:51:43Z
modified: 2026-11-12T08:43:35Z
---
Tuesday 20 March, 14:06. The Board's shutdown and a newer plan both stopped the same planning run, passplan 2's work on one plan; the second `Stop` closed a channel that was already closed, and the planner died with a panic.

`sync.OnceFunc(f)` returns a function that calls `f` the first time it is called. Every later call, from any goroutine, waits until that first call has finished and does nothing more. It makes a stop safe to call twice: closing a channel twice panics, and closing it once never does.

```go
stow := sync.OnceFunc(func() {
	fmt.Println("CWR: dish stowed")
})
stow()
stow() // prints nothing
```

## Your task

`newRun` sets `r.stop` to a function that closes `r.done`. Wrap that function in `sync.OnceFunc`, so `Stop` can be called any number of times, from any goroutine.

> [!NOTE]
> `Done` hands out `r.done` as a `<-chan struct{}`: whoever watches the run can wait on it, and only `Stop` can close it.
