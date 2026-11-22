---
slug: merge-eighty-streams
title: Merge Eighty Streams
type: code_tests
position: 6
difficulty: 9
topics: generics, iterators
created: 2026-11-11T12:51:43Z
modified: 2026-11-12T08:43:35Z
---
Wednesday 11 April. The shadow run's long sorted lists give way to streams, one per satellite and eighty in all, each predicted a day at a time, and passplan 2 must merge them a pass at a time.

To merge sorted sequences, keep each one's next value, its head, in a priority queue. Pop the smallest head, yield it, and pull that sequence's next value into the queue. `iter.Pull` turns a sequence into a `next` function, and a `stop` to call once you're done:

```go
next, stop := iter.Pull(winds)
defer stop()
if kmh, ok := next(); ok {
	fmt.Println("first gust:", kmh)
}
```

## Your task

`MergeBy` yields one sequence after another. Merge them instead:

1. Pull each sequence, `defer` its `stop`, and push its first value as a `head` with the sequence's index into a `Queue` ordered by `byHead(cmp)`.
2. Pop the smallest head and yield its value, until the queue is empty or `yield` returns false.
3. After each yield, push that sequence's next value as its new head.
