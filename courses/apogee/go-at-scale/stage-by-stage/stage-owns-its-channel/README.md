---
slug: stage-owns-its-channel
title: A Stage Owns Its Channel
type: code_tests
position: 0
difficulty: 7
topics: channels, goroutines
created: 2026-11-10T15:39:51Z
modified: 2026-11-10T15:39:51Z
---
Monday 5 March. Twenty stations will mean twenty streams of satellite readings, one per station: telemetry feeds. tlmd, the program that takes them in, reads every feed in one loop, and Nkem wants it rebuilt as a pipeline before her drill on Wednesday.

A pipeline is a chain of stages joined by channels. A stage starts a goroutine and returns the channel it sends on. The stage owns that channel: only its goroutine sends on it, and closes it once nothing is left to send, so the loop reading it can end.

```go
// ends when winds' goroutine closes it
for kmh := range winds(ctx, "CWR") {
	fmt.Println("Cape Wren", kmh, "km/h")
}
fmt.Println("feed over")
```

## Your task

`read` never closes its channel, so a loop over it never ends. Make its goroutine begin with `defer close(out)`. Leave `ctx` for now: on Wednesday, stages learn to stop early.

> [!WARNING]
> Close in the goroutine, not in `read`: `read` returns as soon as its goroutine starts, and a send on a closed channel panics.
