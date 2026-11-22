---
slug: run-until-told-to-stop
title: Run Until Told to Stop
type: code_tests
position: 2
difficulty: 9
topics: channels, context, http, select
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
On launch day the deploy tool stops the old process with a signal, the operating system's way of telling a program to stop, and `main` turns that signal into a canceled context. What happens after it belongs in one function that the tests can call too.

`run` holds the server's whole life: it serves in a goroutine, waits for whichever comes first, `ctx` ending or `Serve` failing, and then stops. `main` makes the context and calls it:

```go
ctx, cancel := signal.NotifyContext(
	context.Background(),
	os.Interrupt, syscall.SIGTERM)
defer cancel()
err := run(ctx, ln, app, timeouts)
```

## Your task

`run` starts serving and returns at once. Call `serve` in a goroutine that sends its result on a channel, then `select` on that channel and `ctx.Done()`: return a `Serve` failure as it is, and when `ctx` ends, return `stop(srv, t.Shutdown)`.
