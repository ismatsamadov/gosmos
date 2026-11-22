---
slug: drain-in-order
title: Drain in Order
type: code
position: 4
difficulty: 8
topics: generics, iterators, methods
created: 2026-11-11T12:51:43Z
modified: 2026-11-12T08:43:35Z
---
Ingrid wants the test Board, the copy of the Board that shows passplan 2's plan, to list Cape Wren's Monday straight from the planner's queue, earliest pass first.

A method can return an iterator. `Drain` returns an `iter.Seq[T]`, the function a `range` loop calls with its `yield`. It pops and yields until the queue is empty, or until `yield` returns false because the loop has stopped:

```go
func (r *Ring[T]) All() iter.Seq[T] {
	return slices.Values(r.items)
}

for kmh := range gusts.All() {
	fmt.Println(kmh, "km/h")
}
```

## Your task

`Drain` yields nothing yet. Make it pop and yield until the queue is empty or `yield` returns false. `main` queues Cape Wren's candidates for Monday 9 April, the passes that climb to 10° or more, and prints the first five:

```text
00:18:07 pipit-2/CWR/21230
00:26:02 cirrine-1/CWR/18668
01:20:50 brine-3/CWR/16710
01:37:36 furrow-a/CWR/23309
02:00:58 cirrine-1/CWR/18669
```

> [!NOTE]
> `main` breaks after five passes. A `Drain` that yields again after that makes the program panic.
