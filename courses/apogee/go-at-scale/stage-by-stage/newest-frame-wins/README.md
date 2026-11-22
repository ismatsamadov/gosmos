---
slug: newest-frame-wins
title: The Newest Frame Wins
type: code_tests
position: 4
difficulty: 7
topics: channels, maps
created: 2026-11-10T15:39:51Z
modified: 2026-11-10T15:39:51Z
---
Ystmark and Tirrin Rock (TIR), a new station on a rock in the Norwegian Sea, will both hear Cirrine-1. In the drill, its frame 9013 came through Ystmark before frame 9012 came through Tirrin Rock's slower link, and the Board showed 9012 as the newest.

The stage at the end of a pipeline is a sink: it ranges over its input until the input closes, then returns what it gathered. Frames from two stations arrive in any order, so arrival can't tell you which frame is newest. The sequence number, the counter in every frame, can.

```go
func heardBy(in <-chan Heard) map[string]int {
	n := map[string]int{}
	for h := range in {
		n[h.Station]++
	}
	return n
}
```

## Your task

`newest` keeps the last frame it receives for each satellite. Make it keep the frame with the highest `Seq` instead.

> [!TIP]
> A satellite's first frame has nothing to beat, even when its number is 0.
