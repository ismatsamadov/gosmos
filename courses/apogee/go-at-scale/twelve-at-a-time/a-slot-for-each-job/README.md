---
slug: a-slot-for-each-job
title: A Slot for Each Job
type: code_tests
position: 0
difficulty: 8
topics: channels, context, goroutines, select
created: 2026-11-11T12:51:42Z
modified: 2026-11-12T08:43:35Z
---
Nkem starts you on passplan 2's predictions: a goroutine for each of the eighty satellites is fine, as long as only twelve of them predict at once, one in each planner slot.

A buffered channel makes a semaphore, a counter of free slots. Its room is the number of slots: a send takes one, and waits while every slot is taken; a receive gives one back.

```go
dishes := make(chan struct{}, 2) // Ystmark has 2
for _, sat := range sats {
	dishes <- struct{}{} // waits while both track
	go func() {
		defer func() { <-dishes }() // free again
		track(sat)
	}()
}
```

## Your task

`Acquire` and `Release` do nothing, so the given `predictAll` starts all eighty predictions at once. Make `Acquire` send into `s` in a `select` beside `<-ctx.Done()`, so a waiting job gives up with `ctx.Err()` when its run is canceled, and make `Release` receive from `s`.

> [!TIP]
> A `struct{}` holds nothing and takes no memory: the channel's room is all that counts.
