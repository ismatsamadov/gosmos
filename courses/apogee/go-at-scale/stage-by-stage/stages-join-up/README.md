---
slug: stages-join-up
title: Stages Join Up
type: code
position: 1
difficulty: 7
topics: channels, goroutines
created: 2026-11-10T15:39:51Z
modified: 2026-11-10T15:39:51Z
---
Wednesday 7 March, 21:00. Nkem starts the drill, her rehearsal of twenty feeds at once, with a single feed: Port Ellery's day, five lines. Fog garbled one, and one shows Pipit-1 in safe mode, its minimal state, where Ingrid put it for the drill.

Stages join up when one stage's channel becomes the next stage's input. Each runs in its own goroutine, so one line can be parsed while the next is still being read, and each channel closes after the one that feeds it.

```go
gusts := winds(ctx, "CWR")
alarms := over(ctx, gusts, 72) // reads gusts
for kmh := range alarms {
	fmt.Println("stow at", kmh, "km/h")
}
```

## Your task

`main` prints the raw lines `read` sends. Pass `read`'s channel to `parse` instead, and print each frame `parse` sends, like this:

```text
pipit-1 6203 NOMINAL 7.94 V
pipit-1 6204 SAFE 7.93 V
cirrine-1 9003 NOMINAL 8.04 V
cirrine-1 9005 NOMINAL 8.03 V
```

> [!NOTE]
> `parse` drops the line it can't read: the fog turned a 0 in frame 9004 into the letter O.
