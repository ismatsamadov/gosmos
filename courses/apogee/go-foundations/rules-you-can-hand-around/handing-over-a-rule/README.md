---
slug: handing-over-a-rule
title: Handing Over a Rule
type: code_tests
position: 1
difficulty: 3
topics: functions
created: 2026-11-04T16:56:59Z
modified: 2026-11-05T08:05:37Z
---
Ingrid's checklist works for Furrowlight, but Cirrine's satellite needs a different minimum angle. Nkem wants the caller to choose the rule, rather than another customer name buried in the checklist.

A function parameter can itself be a function. Its type tells you how to call it; the caller supplies the behavior. In this example, `inspect` calls the supplied `accept` function with an integer and returns its boolean result:

```go
func inspect(n int, accept func(int) bool) bool {
	return accept(n)
}
```

## Your task

Complete `checkPass`: call its `rule` parameter with `elevationDeg` and return the result. The tests supply different rules, including one that accepts low passes, so do not write a minimum angle into `checkPass`.
