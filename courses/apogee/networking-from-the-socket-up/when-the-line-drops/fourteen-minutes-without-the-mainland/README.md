---
slug: fourteen-minutes-without-the-mainland
title: Fourteen Minutes Without the Mainland
type: code_tests
position: 8
difficulty: 9
topics: context, errors, networking, time
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Friday 12 May, 10:00. Tomás's queue for Dish 2 is running, out of stow, point, a 0.35° nudge called a jog, stow again, when Sunniva pulls the island's line mid-jog for fourteen minutes, as Tuesday's storm did.

`Drain` sends a queue in order through `Send`, one outcome per command. A refusal, a `*NakError`, is an outcome, and the queue goes on. Any other error, such as the context ending, stops it: `Drain` returns the outcomes so far with that error.

```go
err := l.Send(ctx, unstow)
_, refused := errors.AsType[*NakError](err)
fmt.Println(refused, err)
// true seq 119 refused: wind lock
```

## Your task

`Drain` in `link.go` sends only the first command. Send each command of `queue` in order with `l.Send`, keep each outcome, and stop early only on an error that is not a `*NakError`.

> [!NOTE]
> The test plays the drill on the simulator in a `testing/synctest` bubble, where fourteen minutes of fake time pass in no real time.
