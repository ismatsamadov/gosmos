---
slug: never-wait-for-the-board
title: Never Wait for the Board
type: code_tests
position: 4
difficulty: 8
topics: channels, select
created: 2026-11-11T12:51:42Z
modified: 2026-11-11T12:51:42Z
---
Wednesday afternoon. One of the Board's screens hung for a minute, its queue of four schedules filled up, and passplan 2 stopped to wait for it. Ingrid: "The Board needs the newest plan. Drop old plans, never the planner."

A `select` with a `default` case never waits: when no other case is ready, `default` runs at once. A send that must not wait tries the queue that way; when it is full, it takes the oldest value out, also without waiting, and tries again.

```go
select {
case f := <-frames:
	fmt.Println("next frame", f.Seq)
default:
	fmt.Println("nothing waiting") // no wait
}
```

## Your task

`offer` sends with a plain `q <- s`, which waits while the queue is full. Try the send in a `select` with `default` instead; when it can't go, take one schedule out of `q` without waiting, count it as dropped, and try again. Return the count.

> [!WARNING]
> `default` alone would drop the new schedule and keep the old ones, and the Board would stay on an old plan.
