---
slug: quick-check-the-weeks-plan
title: Quick Check: The Week's Plan
type: quiz
position: 3
difficulty: 8
topics: maps, slices, time
created: 2026-11-10T20:36:24Z
modified: 2026-11-11T06:44:13Z
---
On Wednesday the rehearsal starts calling the stations. First, Nkem has three questions about the plan.

This is the heart of `assign`, in `plan.go`, for one pass `p`:

```go
isFree := func(t time.Time) bool {
	return !p.AOS.Before(t)
}
d := slices.IndexFunc(next, isFree)
if d < 0 {
	dropped++ // no dish is free
}
```

## Your task

Answer three questions: why a high pass can still be dropped, why the stations can be planned in any order, and why a booking's AOS must match to the second.
