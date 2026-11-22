---
slug: half-the-garbage
title: Half the Garbage
type: code_tests
position: 9
difficulty: 9
topics: maps, performance, slices
created: 2026-11-12T03:10:33Z
modified: 2026-11-13T16:25:47Z
---
Friday 25 May. The third wave goes live in ten days, and Nkem's last target is passplan 2's per-station step, `planStationFat`: on Cape Wren's day it makes 1,471 allocations and 213 KB of garbage. She wants half of both.

Three of this week's cuts have a place in it. Its copy of the candidates grows from nil, where a scratch from the pool is ready. Every key is a formatted string, where a `passKey` allocates nothing. Its map of passes seen grows as it fills, and a map takes a size too:

```go
byCode := make(map[string]Site, len(n.Sites))
for _, s := range n.Sites {
	byCode[s.Code] = s
}
fmt.Println(byCode["CWR"].Dishes) // 1
```

## Your task

`planStationLean` is a copy of `planStationFat`. Gather its candidates in a scratch from `getScratch`, put it back with `putScratch` when you're done, make its map with room for every candidate, and key the map by `passKey`. Its contacts must stay exactly `planStationFat`'s.

> [!NOTE]
> The tests measure both on Cape Wren's and Ystmark's days with `testing.Benchmark` and want half the allocations and half the bytes.
