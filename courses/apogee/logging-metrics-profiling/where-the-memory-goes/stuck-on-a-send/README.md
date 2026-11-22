---
slug: stuck-on-a-send
title: Stuck on a Send
type: code_tests
position: 6
difficulty: 10
topics: channels, defer, goroutines
created: 2026-11-10T15:36:39Z
modified: 2026-11-10T15:36:39Z
---
Thursday 8 February. Four of the goroutines you saved at 02:00 wait on one line of `runPass`, `done <- result{rep, err}`, each holding one of the station's four workers, the slots uplinkd flies its passes on. With all four held, every later pass waited for one until LOS.

A send on a channel with no room waits for a receiver. `runPass` stops receiving at LOS, so when the radio gives up later, the goroutine waits forever and never reaches the line that gives its worker back. A channel with room for one takes the send with nobody receiving, and `defer` gives the worker back however the goroutine ends:

```go
acks := make(chan int, 1) // room for one
acks <- 4688              // no receiver needed
fmt.Println(len(acks))    // 1
```

## Your task

In `runPass`, give `done` room for one result, and give the worker back with a `defer` placed right after the goroutine takes it, instead of at the goroutine's end.

> [!NOTE]
> The test counts what 20 failed passes leave behind with your `leftBehind`, in `memory.go`.
