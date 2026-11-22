---
slug: a-rule-in-a-variable
title: A Rule in a Variable
type: code
position: 0
difficulty: 3
topics: functions, types, variables
created: 2026-11-04T16:56:59Z
modified: 2026-11-04T16:56:59Z
---
The week after Pipit-3's launch, Ingrid wants Furrowlight's rule on the pass checklist: the satellite must rise at least 20 degrees above the horizon. The rule is written; now the checklist needs to keep it.

A function is a value. Assign its name without parentheses to keep the function; add parentheses and arguments to call it. Its type describes its parameters and results: `func(float64) bool`.

```go
report := fmt.Println
report("Ready for the next pass")
```

## Your task

`rule` already holds `highEnough`. Replace the two placeholder `false` values with calls to `rule`, using 15 and 20. Keep `highEnough` unchanged, and remove the unused-variable placeholder.
