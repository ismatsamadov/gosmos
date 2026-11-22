---
slug: a-channel-for-frames
title: A Channel for Frames
type: code
position: 1
difficulty: 6
topics: channels, goroutines
created: 2026-11-01T07:36:46Z
modified: 2026-11-01T07:36:46Z
---
At Ystmark, the poller runs in its own goroutine and hands Pipit-1's frames to the archiver. Sunniva's log shows three frames heard and one archived.

A channel carries values from one goroutine to another: `ch <- v` sends `v`, and `v := <-ch` receives one value. On a channel made with `make(chan T)`, a send waits until another goroutine receives it.

```go
winds := make(chan int)
go func() { winds <- 64 }()
kmh := <-winds
fmt.Println("Ystmark wind", kmh, "km/h")
```

## Your task

`main` receives one sequence number and stops. Receive all three that `poller` sends, printing `archived` and each number.

> [!TIP]
> `make(chan uint32, 8)` gives a channel room for 8 values: its sends wait only when it is full.
