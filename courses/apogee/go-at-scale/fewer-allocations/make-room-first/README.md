---
slug: make-room-first
title: Make Room First
type: code_tests
position: 0
difficulty: 8
topics: iterators, performance, slices
created: 2026-11-12T03:10:33Z
modified: 2026-11-12T09:50:21Z
---
Nkem's heap profile of a replan lists `collect`, which gathers each station's candidates for the planner: for Cape Wren's 360 on take-over day, it allocates ten arrays to hold them.

When `append` finds a slice full, it allocates a bigger array, copies every value into it and leaves the old array as garbage. When you know how many values are coming, make the room first: `make([]T, 0, n)` has length 0 and room for `n`, so appends fill it without allocating.

```go
codes := make([]string, 0, len(sites))
for _, s := range sites {
	codes = append(codes, s.Code)
}
fmt.Println(len(codes), cap(codes)) // 20 20
```

## Your task

passplan 2 replans every ten minutes, and a station's candidates rarely change between runs, so `collect` gets `n`, the station's count from the last run. Start its slice with room for `n` passes: one allocation for Cape Wren's day instead of ten.

> [!TIP]
> `n` is a guess. When more passes come, `append` grows the slice as before, and every pass is kept.
