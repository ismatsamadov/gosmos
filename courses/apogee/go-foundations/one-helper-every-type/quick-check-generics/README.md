---
slug: quick-check-generics
title: Quick Check: Generics
type: quiz
position: 2
difficulty: 5
topics: generics
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:27Z
---
Before generics reach the scheduler, Nkem asks two questions. The first is about a helper someone pushed this morning, the second about a call to your `Filter`:

```go
func Bigger[T any](a, b T) T {
	if a > b {
		return a
	}
	return b
}

high := Filter(passes, aboveTen)
```

`passes` is a `[]Pass`, and `aboveTen` is a `func(Pass) bool` that keeps the passes rising above 10 degrees.
