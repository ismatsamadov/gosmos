---
slug: the-command-in-flight
title: The Command in Flight
type: code_tests
position: 5
difficulty: 9
topics: networking, pointers
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
At 07:22:04 on Tuesday, a 0.35° nudge of Dish 2, Tomás's calibration jog, left dishd; the line dropped two seconds later. The restarted dishd sent it again under a new number, and Dish 2 moved twice.

A command sent but unanswered is in flight: `l.inflight`. After a reconnect, the WELCOME settles it: only dishd sends commands, one at a time, so a `LastSeq` at or past the command's number means it ran. If not, it goes again under its own number, which the controller would refuse had it run:

```go
// the controller's check
if seq <= lastSeq {
	return nak(1, "duplicate") // nothing moves
}
```

## Your task

`resume` in `link.go` always returns `l.inflight`. When `w.LastSeq` is at or past `l.inflight.Seq`, set `l.inflight` to nil and return nil; otherwise return it unchanged.

> [!NOTE]
> The test drops the simulator's line once after the jog runs and once before: either way, Dish 2 must jog once.
