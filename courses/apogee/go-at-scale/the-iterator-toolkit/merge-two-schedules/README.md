---
slug: merge-two-schedules
title: Merge Two Schedules
type: code_tests
position: 6
difficulty: 9
topics: generics, iterators, sorting
created: 2026-11-11T22:52:57Z
modified: 2026-11-12T08:43:35Z
---
Wednesday afternoon. Nkem wants the shadow run's log to list both planners' contacts as one stream in AOS order, so that each contact the two agree on shows twice in a row, the old planner's copy first.

To merge two sorted sequences, pull from both and keep each one's head, the value pulled and not yet yielded. Yield the smaller head and pull the next from its sequence; when one ends, yield the rest of the other. On a tie, yielding `a`'s head first keeps equal values in the order of their sequences: the merge is stable. A head is a value and its `ok`:

```go
next, stop := iter.Pull(slices.Values(kmh))
defer stop()
g, ok := next() // the head
for ok && g < 72 {
	g, ok = next()
}
fmt.Println("first stow gust:", g, ok)
```

## Your task

`Merge` yields all of `a`, then all of `b`. Make it pull from both and yield their values in the order its comparison function sets, `a`'s first on a tie, with both stopped on every way out.
