---
slug: a-heap-of-passes
title: A Heap of Passes
type: code_tests
position: 1
difficulty: 8
topics: interfaces, methods, slices
created: 2026-11-11T12:51:43Z
modified: 2026-11-11T12:51:43Z
---
The shadow run sorts a hundred days of candidates, the passes that climb high enough to be worth a dish, before it plans the first one. Nkem, at the whiteboard: "A queue that always knows its next pass, without sorting a hundred days."

`container/heap` keeps a slice in heap order: the smallest value sits at index 0, and a push or a pop moves only a few values. Your type gives it `Len`, `Less` and `Swap`, plus `Push` and `Pop` for the slice's end. You call `heap.Push` and `heap.Pop`, and they call your methods:

```go
h := &gusts{72, 50, 81}
heap.Init(h)
heap.Push(h, 64)
fmt.Println(heap.Pop(h)) // 50
```

## Your task

`passHeap` orders passes with `byAOS`. Write its `Push`, which appends `x` as a `PassView`, and its `Pop`, which removes the last pass and returns it.

> [!WARNING]
> Remove the last pass, not the first: `heap.Pop` has already swapped the earliest pass to the end.
