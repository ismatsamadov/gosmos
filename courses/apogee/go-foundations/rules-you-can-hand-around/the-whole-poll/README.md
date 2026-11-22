---
slug: the-whole-poll
title: The Whole Poll
type: code_tests
position: 7
difficulty: 4
topics: functions, loops
created: 2026-11-04T16:56:59Z
modified: 2026-11-04T16:56:59Z
---
Ingrid's evening checklist must satisfy every customer's rule before she calls GO. Nkem wants it to stop at the first refusal, so later checks do not run after a pass is rejected.

`type Rule func(float64) bool` gives the function type a name. The supplied `[]Rule` is a slice, an ordered list of rules. `for _, rule := range rules` visits each function in order; the underscore discards its position. You will build and grow slices in the schedule chapter.

```go
type Reading func() int
read := Reading(func() int { return 7 })
fmt.Println(read())
```

## Your task

In `poll`'s supplied loop, call each rule with `elevationDeg`. Return `false` immediately if it refuses. Keep the final `true`: an empty checklist has no refusal.
