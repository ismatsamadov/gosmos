---
slug: a-week-exactly
title: A Week, Exactly
type: code_tests
position: 7
difficulty: 9
topics: loops, metrics
created: 2026-11-10T15:36:37Z
modified: 2026-11-21T12:23:38Z
---
Tuesday 19 December, 10:00. The pager's 2,017 readings of Lantern Atoll's counters for the week are in, and Yusra waits for one number to send Brinewatch. On Monday at 08:00 the Monday restart, a line on Ingrid's checklist since November, restarted the station's uplinkd, and every counter began again at 0.

A counter starts again at 0 when its program restarts. Over a week, the total is the sum of the increases from each reading to the next, and a reading lower than the one before means a restart: the new value is what the program has counted since. Last minus first misses everything before the restart, and with `uint64`, going below 0 wraps around:

```go
var first, last uint64 = 3, 1
fmt.Println(last - first)
// 18446744073709551614
```

## Your task

`increase` returns the last reading minus the first. Add up the increase from each reading to the next instead, counting a drop as a restart: readings 3, 5, 0, 1 make 2 + 0 + 1, which is 3.

> [!NOTE]
> Fewer than two readings make 0: there is no increase to measure.
