---
slug: one-satellite-two-stations
title: One Satellite, Two Stations
type: code_tests
position: 1
difficulty: 9
topics: channels, mutexes, select
created: 2026-11-12T02:35:46Z
modified: 2026-11-12T11:16:23Z
---
Tuesday afternoon. Nkem: "What does the satellite see if two stations send it a command at once?" Over the pole, Ystmark and Tirrin Rock often hold the same satellite.

Two passes numbering from their own WELCOMEs can send one number twice: the satellite refuses the second as a replay. A channel with room for one is a lock a `select` can wait on:

```go
mic := make(chan struct{}, 1)
select {
case mic <- struct{}{}:
	speak()
	<-mic
case msg := <-inbox:
	read(msg)
}
```

## Your task

Keep each satellite to one command in flight, across stations: one lock per satellite, awaited while the pass reads its call, let go once its commands are answered, and at LOS. Number each command above both the WELCOME's `last` and the last number sent to its satellite.

The tests fly Thursday at Ystmark and Tirrin Rock, then the week: no replays, a PING from every welcomed pass, and no satellite answering twice within 2 s, the time one answer takes. A line dead until LOS mustn't keep its lock.

> [!WARNING]
> A `sync.Mutex` held across an answer stops the tests' simulated clock: a goroutine waiting on one never counts as blocked.
