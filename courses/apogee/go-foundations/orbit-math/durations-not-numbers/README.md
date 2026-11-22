---
slug: durations-not-numbers
title: Durations, Not Numbers
type: code
position: 7
difficulty: 4
topics: conversions, time
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
A tool read Pipit-1's 552-second pass as 552 minutes and booked Port Ellery for nine hours. Nkem's new rule: a length of time is a `time.Duration`, never a bare number.

A `time.Duration` counts nanoseconds. Multiply by a unit such as `time.Second` and the unit travels with the number; `%v` prints it readably. An `int` variable needs converting first, because `*` wants the same type on both sides:

```go
fmt.Printf("%v\n", 90*time.Second) // 1m30s
n := 3
d := time.Duration(n) * time.Hour
fmt.Printf("%v\n", d) // 3h0m0s
```

## Your task

Fix `passLength` so it returns `minutes` minutes plus `seconds` seconds. `main` then prints `9m12s`, `45s` and `12m0s`.
