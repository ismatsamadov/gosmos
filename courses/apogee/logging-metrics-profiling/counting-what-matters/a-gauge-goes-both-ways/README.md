---
slug: a-gauge-goes-both-ways
title: A Gauge Goes Both Ways
type: code_tests
position: 1
difficulty: 8
topics: defer, metrics
created: 2026-11-10T15:36:37Z
modified: 2026-11-10T15:36:37Z
---
Ingrid wants `passes_running` on the Board: the passes under way at each station right now. On the lab bench it read 1 after a failed pass, and stayed at 1.

A reading that goes up and down, such as passes running or how full a disk is, is a gauge. Whatever adds to a gauge must take its share away on every way out of the function, and `defer` sees to that: a deferred call runs when the function returns, through whichever `return`, and even when it panics:

```go
func fly(ok bool) error {
	defer fmt.Println("worker back")
	if !ok {
		return errors.New("no ack")
	}
	return nil
}
// fly(false) prints "worker back" too
```

## Your task

`track` adds 1 to `g` before `run` and takes it away only when `run` succeeds. Take it away with `defer` instead, so `g` comes back down after a failure and even after a panic, and return `run`'s error unchanged.

> [!NOTE]
> `Gauge`, in `metrics.go`, keeps its reading in an `atomic.Int64`, so 20 passes can move it at once.
