---
slug: twenty-feeds-one-stream
title: Twenty Feeds, One Stream
type: code_tests
position: 5
difficulty: 8
topics: channels, context, goroutines, select
created: 2026-11-10T15:39:51Z
modified: 2026-11-11T12:51:41Z
---
The drill sends twenty feeds into one parse stage: five from the stations Apogee has, fifteen made from where the new ones will stand. Read in turn, as the old loop reads them, one quiet feed holds up the other nineteen.

A fan-in forwards each input from a goroutine of its own, and closes the output from one more goroutine once every forwarder has returned. Each forwarder selects on `ctx.Done()` for its receive and its send, so a cancel ends them all, then the closer.

```go
var wg sync.WaitGroup
for _, code := range []string{"PEL", "TIR"} {
	wg.Go(func() { poll(ctx, code) })
}
wg.Wait() // every poll has returned
```

## Your task

`merge` returns a channel nobody writes to. Forward each feed into it from a goroutine started with `wg.Go`, selecting on `ctx.Done()` as `parse` does, and close it from one more goroutine once `wg.Wait()` returns.

> [!WARNING]
> Return `out` at once: a `merge` that waits for its feeds before returning has nobody reading `out`, and stalls as soon as it fills.
