---
slug: both-dishes-one-service
title: Both Dishes, One Service
type: code_tests
position: 7
difficulty: 9
topics: channels, context, goroutines, networking, select
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Friday 28 April, ten to midnight, the sun still up over the ridge: dishd runs Dish 1 and Dish 2 together for the first time, and Sunniva watches both from the hut.

`Run` puts the week together. It starts the watch port, then a `select` waits for whichever comes first: the next tick or the end of the context. On a tick, `pollAll` asks both dishes at once and `Broadcast` sends each dish's `watchLine`; when the context ends, `Shutdown` closes everything before `Run` returns:

```go
for {
	select {
	case <-ctx.Done():
		return nil
	case now := <-ticks:
		log.Println("tick", now)
	}
}
```

## Your task

`Run` returns `nil` at once. Start the watch port with `s.Start(s.ln)`, then loop: on each tick, poll every dish and broadcast each result's `watchLine`, Dish 1 first; when `ctx` ends, call `s.Shutdown()` and return `nil`.

> [!WARNING]
> Returning without `Shutdown` leaves every watcher and both controllers connected to a dishd that has stopped.
