---
slug: a-ring-of-readings
title: A Ring of Readings
type: code_tests
position: 5
difficulty: 7
topics: generics, methods
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:27Z
---
The Fishbowl's alert panel shows the last few of everything: Pipit-1's voltages, Thistle's frames, the latest pass IDs. Nkem wants one ring buffer to hold them all.

A type can have type parameters too. Its methods name them again on the receiver, without the constraint:

```go
// Pair holds two values of one type.
type Pair[T any] struct{ A, B T }

// Swap returns the values the other way round.
func (p Pair[T]) Swap() (T, T) { return p.B, p.A }

p := Pair[string]{"AOS", "LOS"}
fmt.Println(p.Swap()) // LOS AOS
```

## Your task

Implement `NewRing`, `Push` and `Items` in `main.go`. Until the ring is full, `Push` appends. After that, it overwrites the oldest item, at `next`, and moves `next` on. `Items` returns the items in a new slice, oldest first.
