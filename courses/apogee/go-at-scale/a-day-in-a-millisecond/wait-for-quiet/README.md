---
slug: wait-for-quiet
title: Wait for Quiet
type: code_tests
position: 5
difficulty: 9
topics: mutexes, time
created: 2026-11-11T16:51:57Z
modified: 2026-11-12T08:43:36Z
---
Friday 27 April. In a gale at Ostrel Island (OST), a new station in the southern Indian Ocean, the wind hovered at the speed where the dish parks for safety, and the station sent eleven stow and unstow messages in four seconds. passplan 2 replanned eleven times.

A debouncer waits for quiet: each event starts its wait over, and only when none has come for a while does it act, once. `time.AfterFunc(d, f)` runs `f` in its own goroutine after `d`. On that timer, `Reset(d)` starts the wait over, even after `f` has run, and `Stop` cancels a call still waiting.

A replan that waits for quiet:

```go
t := time.AfterFunc(5*time.Second, replan)
t.Reset(5 * time.Second) // 5 s from now
t.Stop()                 // replan never runs
```

## Your task

`Notify` calls `fire` at once. Holding `d.mu`, make the first `Notify` start a `time.AfterFunc` of `d.quiet` that calls `d.fire`, and each later one `Reset` it. `Stop` stops the timer, if there is one.

> [!NOTE]
> Messages come from every station's goroutine, hence the mutex: two first calls of `Notify` must never both make a timer.
