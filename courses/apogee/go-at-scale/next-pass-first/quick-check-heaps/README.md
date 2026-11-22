---
slug: quick-check-heaps
title: Quick Check: Heaps
type: quiz
position: 2
difficulty: 8
topics: interfaces, slices
created: 2026-11-11T12:51:43Z
modified: 2026-11-11T12:51:43Z
---
Before the queue goes into the planner, Nkem writes two lines on the whiteboard and one question beside them: whose `Pop` is whose?

```go
heap.Push(h, p)
next := heap.Pop(h).(PassView)
```

Here `h` is a `*passHeap`, with the `Push` and `Pop` methods you wrote, and the shadow run's queue will hold about 20,000 passes at once.

## Your task

Answer Nkem's three questions: who calls your `Pop` method, why it takes the last element, and how the work of one pop grows with the queue.
