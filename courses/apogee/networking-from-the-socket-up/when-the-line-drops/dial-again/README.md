---
slug: dial-again
title: Dial Again
type: code_tests
position: 0
difficulty: 8
topics: context, loops, networking
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Ingrid, who runs every pass from the Fishbowl, has Tuesday's log open: dishd, the dish service, dialed Dish 2's controller once at 07:22, failed, and never tried again. She restarted it herself at 07:51, fifteen minutes after the line came back.

To stay connected, dial in a loop: when a dial fails, wait, then dial again, until one works. dishd's `Link` waits with `l.Wait(ctx, d)`, which sleeps for `d`, or returns the context's error as soon as the context ends, the caller's way of saying stop:

```go
ctx, cancel := context.WithCancel(ctx)
cancel() // dishd is shutting down
err := l.Wait(ctx, time.Second)
fmt.Println(err) // context canceled
```

## Your task

`redial` in `link.go` dials once. Make it dial again after `l.Wait(ctx, time.Second)` until a dial works, and when a wait fails, return the wait's error.

> [!NOTE]
> Not the last dial's error: a refused dial says the line is down, `context.Canceled` says stop, and the caller acts on the difference.
