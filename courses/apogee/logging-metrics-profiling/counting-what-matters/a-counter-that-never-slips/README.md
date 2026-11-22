---
slug: a-counter-that-never-slips
title: A Counter That Never Slips
type: code_tests
position: 0
difficulty: 8
topics: goroutines, metrics
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:38Z
---
Yusra can send Brinewatch one number for November, and she has four. Nkem: "Then uplinkd keeps the count, and never loses one."

A number a program keeps about its own work is a metric, and one that only ever goes up, such as passes failed, is a counter. uplinkd adds to its counters from many goroutines at once, and `c.n++` takes three steps, a read, an add and a write: two goroutines can read the same old value, and one pass goes uncounted. An `atomic.Uint64`, from `sync/atomic`, adds in one step that nothing can split:

```go
var frames atomic.Int64
frames.Add(3)
frames.Add(-1)
fmt.Println(frames.Load()) // 2
```

## Your task

Make `Counter`'s field `n` an `atomic.Uint64` instead of a plain `uint64`, then write `Inc`, `Add` and `Value` with its `Add` and `Load` methods.

> [!NOTE]
> A plain `uint64` loses a count only now and then, and a test can pass by luck, so the test also reads `metrics.go` for the `atomic.Uint64`.
