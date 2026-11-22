---
slug: one-goroutine-per-dish
title: One Goroutine per Dish
type: code_tests
position: 0
difficulty: 8
topics: goroutines, networking, slices
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Monday 24 April: with two controllers in, dishd asks Dish 1 for its status, then Dish 2. On a backhaul day, at 0.7 s a question, every screen ends up 1.4 s behind.

Each client has its own connection, so both questions can travel at once: a goroutine per dish, and a `sync.WaitGroup` to wait for them. `wg.Go(f)` runs `f` in a new goroutine and counts it; `wg.Wait()` returns once all have finished:

```go
var wg sync.WaitGroup
for i, name := range names {
	wg.Go(func() { lens[i] = len(name) })
}
wg.Wait() // every length is in place
```

## Your task

`pollAll` asks one dish after the other. Ask each dish in its own goroutine with `wg.Go`, store its `Result` at its index in `results`, and wait for all of them before the check that follows.

> [!NOTE]
> Each goroutine writes only its own element of `results`, so none touches another's memory.
