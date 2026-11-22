---
slug: quick-check-slots-and-pools
title: Quick Check: Slots and Pools
type: quiz
position: 2
difficulty: 7
topics: channels, goroutines
created: 2026-11-11T12:51:42Z
modified: 2026-11-11T12:51:42Z
---
Tuesday afternoon, Nkem writes "Slots or workers?" on the whiteboard, above a loop from an old draft.

With `s := newSlots(12)`, the loop takes a slot for each check and gives it back when the check passes:

```go
for _, e := range els {
	s.Acquire(ctx)
	go func() {
		if check(e) == nil {
			s.Release()
		}
	}()
}
```

## Your task

Answer Nkem's three questions: what the thirteenth `Acquire` does while twelve checks run, what a failing check does to its slot, and how many goroutines slots and a pool of twelve workers start for eighty jobs.
