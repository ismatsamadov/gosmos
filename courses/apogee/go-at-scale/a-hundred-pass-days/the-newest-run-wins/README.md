---
slug: the-newest-run-wins
title: The Newest Run Wins
type: code_tests
position: 5
difficulty: 10
topics: context, goroutines, testing
created: 2026-11-12T03:18:38Z
modified: 2026-11-12T03:18:38Z
---
Tuesday 26 June, in the shadow run. A stow arrived halfway through a tick's run, the tick's run finished second, and its plan, blind to the stow, replaced the stow's on the test Board.

The draft `Run` plans each event in a goroutine and numbers the run when it ends, so the last to finish looks newest. Two rules make the newest news win. A run takes its number as its event arrives, so the publisher's `Offer` refuses whatever an older run hands it late. And news cancels the run in progress, whose plan is out of date already:

```go
stop := func() {}
for q := range queries {
	stop() // a newer query: drop the old search
	ctx, cancel := context.WithCancel(parent)
	stop = cancel
	go search(ctx, q)
}
```

## Your task

In `Run`, take each run's number as its event arrives, cancel the run in progress, and plan the new one with a context of its own. Before `Run` returns, wait for every run it started: when `events` closes, the last run still finishes.

> [!TIP]
> A `sync.WaitGroup`'s `Go` method starts a goroutine the group then waits for.
