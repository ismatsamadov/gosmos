---
slug: a-queue-for-any-type
title: A Queue for Any Type
type: code_tests
position: 3
difficulty: 9
topics: generics, interfaces, methods
created: 2026-11-11T12:51:43Z
modified: 2026-11-12T08:43:35Z
---
The planner will need a queue of passes, a queue of contacts and one more inside every merge. Nkem wants one queue for all of them, written once.

A generic type has a type parameter, as `Queue[T any]` does, and its methods use `T` wherever a value goes. `container/heap` deals in `any`; a generic queue keeps that inside, so its callers push and pop `T`s:

```go
type Ring[T any] struct{ items []T }

func (r *Ring[T]) Add(v T) {
	r.items = append(r.items, v)
}

gusts := &Ring[int]{}
gusts.Add(72)
```

## Your task

`heapOf[T]` implements `heap.Interface` for a slice and a comparison. Write `NewQueue`, which gives the queue a `heapOf` holding the comparison it is given, then `Push`, `Pop` and `Len` on `heap.Push`, `heap.Pop` and the heap's own `Len`. An empty queue's `Pop` returns the zero value and false, as the stub does.

> [!TIP]
> `heap.Pop(q.h)` returns an `any`: `.(T)` gives the value back as a `T`.
