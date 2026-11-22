---
slug: count-the-goroutines
title: Count the Goroutines
type: code_tests
position: 5
difficulty: 9
topics: goroutines, profiling, testing
created: 2026-11-10T15:36:39Z
modified: 2026-11-10T15:36:39Z
---
Wednesday, 16:00. At 01:50 the night shift had read you Lantern Atoll's `goroutines`, the number of goroutines its uplinkd runs: 19, where the other stations read 14 or 15. Nkem: "A number caught it in the night. Which test would have caught it first?"

`runtime.NumGoroutine` returns how many goroutines exist right now. Count before a function runs and after it returns: the difference is what it left running. A goroutine can take a moment to end after its function returns, so look again every millisecond, for up to a second, before you decide:

```go
before := runtime.NumGoroutine()
go time.Sleep(5 * time.Millisecond)
fmt.Println(runtime.NumGoroutine() - before)
time.Sleep(20 * time.Millisecond)
fmt.Println(runtime.NumGoroutine() - before)
// 1, then 0
```

## Your task

`leftBehind` returns 0. Make it count goroutines, run `f`, then check every millisecond, for up to a second, until the count is back where it started, and return how many more there still are.

> [!NOTE]
> When `f` ends goroutines that were running before it, the count falls below where it started: nothing was left behind, so the answer is 0.
