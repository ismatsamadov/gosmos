---
slug: any-number-of-readings
title: Any Number of Readings
type: code_tests
position: 8
difficulty: 4
topics: functions, slices
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:28Z
---
The temperature panel averages Thistle's battery temperatures: three readings on a quiet day, three hundred after a busy one. One function should take them all.

A variadic parameter, `...T`, accepts any number of arguments, none included, and inside the function it is a `[]T`. To pass a slice you already have, put `...` after it:

```go
func count(codes ...string) int {
	return len(codes)
}

count()             // 0
count("YST", "PEL") // 2
all := []string{"YST", "PEL", "SQT"}
count(all...) // 3
```

## Your task

Make `average` return the mean of its readings, and 0 when it gets none.

> [!WARNING]
> With no readings, the sum and the count are both 0, and 0 divided by 0 in `float64` is NaN, not 0.
