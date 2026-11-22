---
slug: how-fast-it-burns
title: How Fast It Burns
type: code_tests
position: 4
difficulty: 8
topics: conversions, metrics
created: 2026-11-10T15:36:37Z
modified: 2026-11-10T15:36:37Z
---
On the night of 27–28 December, Lantern Atoll lost six passes, more than its whole month's budget, and no number showed it while it happened. Ingrid wants one for the pager, which wakes whoever is on call: how fast the budget is going.

The burn rate divides the share of passes lost by the share the target allows. At 1, a budget lasts exactly its window; at 50, a month's budget is gone in fourteen hours. Dividing floats never panics, but two answers need watching:

```go
allowed := float64(10000-9900) / 10000
fmt.Println(allowed, 0.25/allowed)
// 0.01 25
zero := 0.0
fmt.Println(1/zero, zero/zero)
// +Inf NaN
```

## Your task

Make `burnRate` divide the share lost, `lost/booked`, by the share `t` allows, `(10000 - t)/10000`. Return 0 when nothing was booked or nothing was lost; a loss against a target of 10000 gives `+Inf`.

> [!TIP]
> `+Inf` is the honest answer at 100%: any loss spends a budget of nothing infinitely fast.
