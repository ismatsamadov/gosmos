---
slug: an-error-with-fields
title: An Error with Fields
type: code_tests
position: 6
difficulty: 6
topics: errors, pointers, structs
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
The pager doesn't want a sentence about Thistle's 51 °C battery. It wants the field, the value and the limits, so it can decide how loudly to page.

An error can be any type with an `Error() string` method, so it can carry fields. `errors.AsType[*LimitError](err)` searches `err`, and everything wrapped inside it, for a `*LimitError`, the type in the square brackets, and returns it with `true` if one is there:

```go
type StowError struct{ Kmh float64 }
func (e *StowError) Error() string {
	return "stowed for wind"
}
err := fmt.Errorf("CWR: %w", &StowError{Kmh: 81})
if se, ok := errors.AsType[*StowError](err); ok {
	fmt.Println(se.Kmh) // 81
}
```

## Your task

1. Make `checkTemp` return a `*LimitError` for `temp_c`, with the reading and both limits.
2. Make `pagerDetail` find that error in `err` and return its fields as a line like `temp_c=51.0 outside [-10.0, 45.0]`, or `""` when there is none.

> [!NOTE]
> Most Go code you'll read still uses the older `errors.As(err, &target)`, which fills a variable you declare first.
