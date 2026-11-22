---
slug: many-watchers-at-once
title: Many Watchers at Once
type: code_tests
position: 2
difficulty: 8
topics: goroutines, networking
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Three screens want the watch port: Sunniva's terminal in the hut, Tomás's laptop and Ingrid's screen in the ops room at Port Ellery, the Fishbowl. Today the second to connect waits until the first hangs up.

An accept loop that serves each connection before accepting the next serves one watcher at a time. Started in a goroutine of its own, each handler runs as long as its watcher stays while the loop goes back to `Accept`. Closing the listener is how the loop is told to stop: `Accept` then fails with an error matching `net.ErrClosed`.

```go
ln.Close()
_, err := ln.Accept()
fmt.Println(errors.Is(err, net.ErrClosed)) // true
```

## Your task

`serve` handles each connection before accepting the next. Start each `handle(c)` with `wg.Go` and keep accepting. Return `nil` once `ln` is closed, and any other `Accept` error as it is.

> [!NOTE]
> The caller passes `wg` so it can wait for every handler: the loop can end before they do.
