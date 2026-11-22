---
slug: severity-levels
title: Severity Levels
type: code_tests
position: 7
difficulty: 5
topics: constants, types
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:26Z
---
The pager, the service that alerts the engineer on call, ranks every unexpected reading, or anomaly, as Info, Warning or Critical. When two rules fire on the same frame, it must send an alert, or page, for the worse of the two.

A named type such as `type Severity int` keeps severities apart from plain numbers. In a `const` block, `iota` is 0 on the first line, 1 on the next, and so on, and a line with no value repeats the type and expression above it:

```go
type Shift int

const (
	Day     Shift = iota // 0
	Evening              // 1
	Night                // 2
)
```

`Severity` is an `int` underneath, so `<` and `>` compare its values.

## Your task

1. Number the three constants with `iota`, so that `Info < Warning < Critical`.
2. Complete `worst` so it returns the more severe of `a` and `b`.
