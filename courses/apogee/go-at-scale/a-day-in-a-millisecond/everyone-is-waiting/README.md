---
slug: everyone-is-waiting
title: Everyone Is Waiting
type: code_tests
position: 6
difficulty: 9
topics: testing, time
created: 2026-11-11T16:51:57Z
modified: 2026-11-11T16:51:57Z
---
Monday 30 April. Your debouncer turned Friday's eleven replans into one, and Nkem wants its test written by you, so no later change breaks it unnoticed.

`synctest.Wait()` returns once every other goroutine in the bubble is blocked or done: the moment to look, with no real time spent. Call it after a `time.Sleep`, before you look, so that whatever came due has finished.

A timer's function runs on its own goroutine, so count with an atomic:

```go
var fires atomic.Int32
fire := func() { fires.Add(1) }
time.AfterFunc(time.Second, fire)
time.Sleep(time.Second)
synctest.Wait()           // fire has run
fmt.Println(fires.Load()) // 1
```

## Your task

`checkDebouncer`, in `checks.go`, is empty. In a bubble, make a debouncer with `newD` and a quiet of 5 seconds, and call `Notify` three times, a second apart. Fail `t` unless nothing has fired a millisecond before the quiet ends, exactly one fire has come when it ends, and after one more `Notify` and a `Stop`, no other fire comes.

> [!NOTE]
> Four mutants must fail your check: one fires at the first `Notify`, one never resets, one fires twice and one ignores `Stop`.
