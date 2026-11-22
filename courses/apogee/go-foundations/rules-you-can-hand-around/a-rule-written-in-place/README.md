---
slug: a-rule-written-in-place
title: A Rule Written in Place
type: code
position: 2
difficulty: 3
topics: functions
created: 2026-11-04T16:56:59Z
modified: 2026-11-05T08:05:37Z
---
Furrowlight's next pass is almost ready to book. Nkem wants its one-use rule beside the checklist entry, where a reviewer can read both together.

A function literal writes a function without declaring a name. It still has parameters, a body and results, and it can be assigned or passed like any other function value:

```go
twice := func(n int) int { return n * 2 }
fmt.Println(twice(4)) // 8
```

## Your task

Change only the body of the function literal assigned to `rule`. Accept an elevation of at least 20 degrees, including exactly 20. Keep the existing calls: all three should use the same function value.

> [!NOTE]
> A declared function variable starts at `nil`, its zero value. Calling it before a function is assigned panics: the program stops with a run-time error instead of returning a default result.
