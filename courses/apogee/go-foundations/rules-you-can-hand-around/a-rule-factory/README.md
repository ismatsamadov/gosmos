---
slug: a-rule-factory
title: A Rule Factory
type: code_tests
position: 5
difficulty: 4
topics: functions
created: 2026-11-04T16:56:59Z
modified: 2026-11-05T08:05:37Z
---
Furrowlight wants passes at least 20 degrees above the horizon; Cirrine wants 30. Nkem wants one rule builder, with the customer's minimum kept inside each rule.

A function can return another function; a function that builds functions this way is called a factory. A parameter of the outer call is also a variable the inner function can capture. Each call builds a new closure with its own parameter:

```go
func addOffset(offset int) func(int) int {
	return func(n int) int { return n + offset }
}
```

## Your task

Complete the returned function in `minElevation`. It must accept elevations greater than or equal to the captured `minDeg`. Tests keep two rules at once; building the second must not change the first.
