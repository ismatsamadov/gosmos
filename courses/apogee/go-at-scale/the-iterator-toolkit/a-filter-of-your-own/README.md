---
slug: a-filter-of-your-own
title: A Filter of Your Own
type: code_tests
position: 1
difficulty: 8
topics: generics, iterators
created: 2026-11-11T22:52:57Z
modified: 2026-11-11T22:52:57Z
---
Monday 2 April. Tirrin Rock's panel on the Board wants its contacts and bookingd wants Furrowlight's, each from a loop of its own with an `if` inside. Nkem: "One idea, written twice."

An adapter takes a sequence and returns another. A filter yields only the values a function says to keep, and returns as soon as its `yield` returns false; that ends its loop over the source, which stops the source too. So a loop that breaks after two values reads no further than it must. A type parameter lets one function serve contacts, numbers and anything else:

```go
func Repeat[V any](v V, n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		for range n {
			if !yield(v) {
				return
			}
		}
	}
}
```

## Your task

`Filter` hands back `seq` unchanged. Make it yield only the values of `seq` that `keep` returns true for, in order, and return as soon as `yield` returns false.

> [!NOTE]
> The tests filter Monday's schedule as passplan 2 plans it: 260 contacts at ten stations, 48 of them Tirrin Rock's.
