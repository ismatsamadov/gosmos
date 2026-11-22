---
slug: keep-the-line-warm
title: Keep the Line Warm
type: code_tests
position: 5
difficulty: 8
topics: channels, networking, time
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Thursday. Sunniva: "Wind 50, steady. Dead line looks like quiet line. Until next command. Maybe an hour." Worse, the controller hangs up on any client that sends nothing for 30 seconds.

A heartbeat settles both: on a timer, dishd sends a PING and requires the controller's PONG within its deadline, proof that the line works both ways. Ticks arrive on a channel, as from a `time.Ticker`, and `range` over a channel ends when it is closed (`Stop` never closes a ticker's):

```go
tk := time.NewTicker(10 * time.Second)
defer tk.Stop()
for now := range tk.C {
	fmt.Println("tick", now.Format("15:04:05"))
}
```

## Your task

`heartbeat` returns `nil` at once. On every tick, send a PING and read the answer, which must be a PONG. Return the first error, or one of your own for any other answer, and `nil` once `ticks` closes.

> [!WARNING]
> Read the PONG every time: one left unread is what the next command reads as its answer.
