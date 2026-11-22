---
slug: map-and-take
title: Map and Take
type: code_tests
position: 2
difficulty: 8
topics: generics, iterators
created: 2026-11-11T22:52:57Z
modified: 2026-11-12T08:43:35Z
---
Monday afternoon. Ingrid's corner of the Board shows the next twelve contacts as lines of text, and the code behind it formats every contact of the day to show twelve.

Two more adapters cover most loops. `Map` changes each value as it passes, often into another type, so it has two type parameters. `Take` yields the first `n` values and then returns, which stops its source too: ask for twelve, and the source makes twelve, not thirteen. Keep any count inside the function an adapter returns, so that each loop over it starts afresh:

```go
func Countdown(from int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for n := from; n > 0; n-- {
			if !yield(n) {
				return
			}
		}
	}
}
```

## Your task

`Map` and `Take` return sequences that yield nothing. Make `Map` yield `f(v)` for each value `v` of `seq`, and `Take` the first `n` values of `seq`, or none at all, without asking `seq`, when `n` is 0 or less. Both return as soon as `yield` returns false.
