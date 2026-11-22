---
slug: the-highest-reading
title: The Highest Reading
type: code_tests
position: 1
difficulty: 6
topics: generics
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:27Z
---
Nkem counts three more loops that differ only in their types: tonight's highest elevation, the warmest battery, the station code that sorts last.

With `any`, Go won't let you compare two `T` values with `<`, because some types have no order. A tighter constraint fixes that: `cmp.Ordered` admits exactly the types that `<` and `>` work on, which are integers, floats and strings.

```go
// Earlier returns a or b, whichever sorts first.
func Earlier[T cmp.Ordered](a, b T) T {
	if b < a {
		return b
	}
	return a
}
```

## Your task

Write `Max` in `main.go`: return the largest item and `true`. For an empty slice, keep returning the zero value and `false`.
