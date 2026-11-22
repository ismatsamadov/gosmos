---
slug: the-whole-board-running
title: The Whole Board, Running
type: code_tests
position: 2
difficulty: 9
topics: context, goroutines, http
created: 2026-11-07T06:23:24Z
modified: 2026-11-07T06:23:24Z
---
Thursday evening: Nkem wants the Board started exactly as the night shift will start it. `main.go` now calls `run`, with a tick every 15 seconds.

`run` owns the program's life. It serves in a goroutine, refreshes when it starts and on every tick, and when its context ends it shuts the server down. `Serve` returns `http.ErrServerClosed` after `Shutdown`: the normal end, never a failure. A refresh that fails keeps the last view, and the night goes on.

## Your task

`run` only serves. Make it:

1. Serve `newServer(b)` on `ln` in a goroutine.
2. Refresh once, then once on every tick.
3. When `ctx` is done, shut down with a 2-second deadline and return nil, or the error `Shutdown` returns. If `Serve` fails before that, return its error.

The screens then read `bye`, and `run` returns within a second.
