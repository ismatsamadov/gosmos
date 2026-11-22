---
slug: one-filter
title: One Filter for Everything
type: code_tests
position: 0
difficulty: 6
topics: functions, generics
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:27Z
---
Nkem finds `filterPasses` in `schedule` and `filterFrames` in `telemetry`: the same eight lines, twice, with only the type changed.

A *type parameter*, in square brackets after a function's name, stands for a type the caller chooses. Its constraint, here `any`, says which types may stand in. Inside the function, `T` goes wherever the type would, and a call rarely names it, as in `LastN[string]`, because Go infers it from the arguments:

```go
// LastN returns the last n items, or all of them
// if there are fewer.
func LastN[T any](items []T, n int) []T {
	return items[max(len(items)-n, 0):]
}

LastN(frames, 5)          // T is TelemetryFrame
LastN([]string{"YST"}, 2) // T is string
```

## Your task

Write `Filter` in `main.go`: return a new slice of the items for which `keep` returns true, in their original order, and leave `items` unchanged. The tests call it with passes, ints and strings.
