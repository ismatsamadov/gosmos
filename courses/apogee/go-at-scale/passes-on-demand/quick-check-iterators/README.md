---
slug: quick-check-iterators
title: Quick Check: Iterators
type: quiz
position: 2
difficulty: 7
topics: iterators
created: 2026-11-11T21:35:30Z
modified: 2026-11-12T08:43:35Z
---
Before the catalog makes its passes on demand, Nkem writes an iterator on the whiteboard and ranges over it in several loops; some stop early.

```go
func countdown(yield func(int) bool) {
	fmt.Println("go")
	for n := 3; n > 0; n-- {
		if !yield(n) {
			return
		}
	}
}
```

Each loop is `for n := range countdown`, and its body prints `n`.

## Your task

Answer Nkem's three questions: when the call `yield(3)` returns, what happens when a loop breaks at 2, and what a second loop over `countdown` prints.
