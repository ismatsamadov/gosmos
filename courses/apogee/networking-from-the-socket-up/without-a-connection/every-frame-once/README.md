---
slug: every-frame-once
title: Every Frame Once
type: code_tests
position: 4
difficulty: 8
topics: methods, networking, structs
created: 2026-11-13T06:15:46Z
modified: 2026-11-13T16:03:24Z
---
Wednesday 14 June, Furrow-A's 08:31 pass. Ingrid's screen showed the battery at 7.86, 7.85, 7.86 volts: on the island's line two datagrams had changed places and one had come twice. Nkem: "What does the screen show if this datagram arrives twice?"

Datagrams can be lost, repeated or overtaken, and nothing says so; each frame's sequence number does. Show a frame only when its number is past the last one shown, and drop the rest: repeats and older frames. A frame that jumps past the next number counts the numbers it skipped as missed, even if one turns up later and is dropped.

```go
a, b := uint32(52818), uint32(52819)
fmt.Println(a-b, b-a) // 4294967295 1
```

## Your task

`Accept` in `tracker.go` shows every frame. Show a frame only when it is the first or its number is past `t.Last`, counting it in `Shown` and the numbers it jumps over in `Missed`; anything else adds one to `Dropped` and returns false.

> [!TIP]
> Unsigned numbers wrap below zero, so subtract the older number from the newer, never the other way.
