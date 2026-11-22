---
slug: counting-without-locks
title: Counting Without Locks
type: code_tests
position: 5
difficulty: 8
topics: goroutines, mutexes
created: 2026-11-11T12:51:42Z
modified: 2026-11-11T12:51:42Z
---
Thursday 22 March. The counts passplan 2 keeps of contacts planned and passes dropped, added from all twelve slots at once, came out different every time Nkem ran the test: they were plain ints.

`n += k` is three steps, a read, an add and a write, so two goroutines can read the same old value and lose a count. An `atomic.Int64`'s `Add` is one step nothing can split, and `Load` never sees half of one. Its zero value is ready to use, and `go vet` objects when a struct holding one is copied.

```go
type Tally struct{ Heard, Lost atomic.Int64 }

var t Tally // ready at zero
t.Heard.Add(12)
fmt.Println(t.Heard.Load(), t.Lost.Load())
// 12 0
```

## Your task

The fields are `atomic.Int64` now, but `Add` and `Snapshot` do nothing yet. Make `Add` add `planned` and `dropped` to their counters and count the run in `Runs`, and make `Snapshot` load all three.

> [!NOTE]
> With plain `int64` fields, `go test -race` in a terminal fails with `WARNING: DATA RACE`. The app runs tests without `-race`, so this test counts instead.
