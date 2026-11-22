---
slug: wait-until-aos
title: Wait Until AOS
type: code_tests
position: 0
difficulty: 8
topics: context, select, time
created: 2026-11-11T16:51:57Z
modified: 2026-11-12T08:43:36Z
---
The station runner, the part of passplan 2 that flies each contact at its AOS, waits by sleeping a minute at a time, so a contact can start up to 59 seconds late.

A `time.Timer` sends on its channel `C` once, when its duration is up; a duration of zero or less sends at once. In a `select` beside `ctx.Done()`, it ends a wait at whichever comes first. `Stop` turns a timer off, so `defer t.Stop()` covers every way out.

This one gives Thistle-2's next frame 90 seconds:

```go
t := time.NewTimer(90 * time.Second)
defer t.Stop()
select {
case f := <-frames:
	fmt.Println("frame", f.Seq)
case <-t.C:
	fmt.Println("no frame for 90 s")
}
```

## Your task

`waitUntil` sleeps in one-minute steps until `now()` passes `at`. Start one timer for `at.Sub(now())` instead, and `select` on its channel and `ctx.Done()`: return `nil` when the timer fires, and `ctx.Err()` when the context ends first. Stop the timer either way.

> [!NOTE]
> The tests run in fake time: their four-hour wait for 04:12:25 takes a millisecond. Tomorrow you write one yourself.
