---
slug: a-set-of-anything
title: A Set of Anything
type: code_tests
position: 6
difficulty: 6
topics: generics, maps, methods
created: 2026-11-01T07:36:46Z
modified: 2026-11-01T07:36:46Z
---
Ingrid asks which stations heard both Thistle and Pipit-1 tonight. The answer should come from a set type that works for anything: station codes today, frame numbers tomorrow.

Map keys must be comparable with `==`, so a type parameter used as a key needs the `comparable` constraint; `any` is refused. A generic type can be built on a map, and its methods use it like any map:

```go
// Tally counts how often each value turns up.
type Tally[T comparable] map[T]int

// Count records one more sighting of v.
func (t Tally[T]) Count(v T) { t[v]++ }
```

## Your task

Implement `Add`, `Has` and `Intersect` in `main.go`. `Intersect` returns a new set of the values in both sets, and leaves both sets as they were.

> [!NOTE]
> `Set` stores `struct{}`, a struct type with no fields, where your sets so far stored `bool`. Its only value, `struct{}{}`, takes no memory.
