---
slug: fresh-input-untimed
title: Fresh Input, Untimed
type: code_tests
position: 3
difficulty: 9
topics: performance, slices, testing
created: 2026-11-11T22:45:16Z
modified: 2026-11-11T22:45:16Z
---
Thursday 10 May. `planDay` sorts its candidates in place, so from the second turn on, passplan 2's benchmark handed it the day it had sorted the turn before: the best case, every time.

Every turn of a benchmark should start from the same input. When the work changes its input, give each turn a fresh copy, and make the copy with the timer stopped: `b.StopTimer()` pauses the clock and the allocation count, and `b.StartTimer()` resumes them, so the copy costs the result nothing:

```go
for b.Loop() {
	b.StopTimer()
	g := slices.Clone(gusts)
	b.StartTimer()
	slices.Sort(g)
}
```

## Your task

`benchFresh` hands `planDay` the same `day` in every turn. In each turn, clone `day` with the timer stopped, start the timer again, and plan the clone.

> [!NOTE]
> `planDay` counts the calls whose candidates came already sorted, and the tests want none. They also want exactly `planDay`'s own allocations per turn.
