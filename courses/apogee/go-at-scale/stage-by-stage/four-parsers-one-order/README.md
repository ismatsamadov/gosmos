---
slug: four-parsers-one-order
title: Four Parsers, One Order
type: code_tests
position: 8
difficulty: 9
topics: channels, context, goroutines
created: 2026-11-10T15:39:51Z
modified: 2026-11-10T15:39:51Z
---
Friday 9 March. Parsing is the drill's slowest stage, and Nkem wants it on four cores, but each satellite's frames must still leave in the order a station heard them.

A worker pool returns results as they finish. To keep the input's order, give each line its own result channel, queued as the line arrives: workers fill them in any order, and one goroutine reads the queue in turn.

```go
replies := []chan string{}
for _, code := range []string{"PEL", "YST"} {
	r := make(chan string, 1)
	replies = append(replies, r)
	go func() { r <- ping(code) }()
}
fmt.Println(<-replies[0], <-replies[1]) // PEL YST
```

## Your task

1. Start `workers` goroutines that take `job`s from one channel and put each `result` in its job's channel.
2. In one goroutine, queue a result channel per line on a channel with room for `workers`, then hand the line to the workers.
3. In another, read the queue in turn and send each frame, skipping refused lines.

> [!TIP]
> Nkem's `send` and `receive`, in `chans.go`, each wrap one `select` beside `ctx.Done()`.
