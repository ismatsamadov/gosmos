---
slug: your-first-function
title: Your First Function
type: code
position: 0
difficulty: 3
topics: functions
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Tomás keeps Apogee's orbit formulas in a battered notebook, and three tools copied his orbits-per-day formula three different ways. He wants it written once.

A function names a piece of work. Its parameters, each with a type, go in parentheses, followed by the type of its result; `return` hands the result back to the caller:

```go
// metersFromKm converts kilometers to meters.
func metersFromKm(km float64) float64 {
	return km * 1000
}
```

Calling `metersFromKm(525)` gives `525000`.

## Your task

1. Where the TODO is, declare `orbitsPerDay(periodMin float64) float64`, returning how many orbits of `periodMin` minutes fit in `minutesPerDay`.
2. In `main`, replace each hard-coded figure with a call to `orbitsPerDay`, passing that satellite's period.
