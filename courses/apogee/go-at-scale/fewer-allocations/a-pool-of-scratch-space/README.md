---
slug: a-pool-of-scratch-space
title: A Pool of Scratch Space
type: code_tests
position: 5
difficulty: 9
topics: goroutines, performance, pointers
created: 2026-11-12T03:10:33Z
modified: 2026-11-12T03:10:33Z
---
Thursday 24 May. Each of the twelve planner slots makes a new 512-pass scratch slice for every station it plans, then drops it. Nkem: "Twelve slots, twelve buffers, not a thousand."

A `sync.Pool` holds values nobody is using, and any goroutine may use it at once. `Get` returns `any`: a value `Put` stored or, when the pool is empty, one that `New` makes. The collector may empty it at any time. Pool pointers to slices: a slice put into an `any` is copied to the heap.

```go
var lines = sync.Pool{New: func() any {
	b := make([]byte, 0, 64)
	return &b
}}
p := lines.Get().(*[]byte)
*p = append((*p)[:0], "AOS 05:29:26"...)
lines.Put(p)
```

## Your task

Keep a `sync.Pool` of `*[]PassView` whose `New` makes the scratch `getScratch` makes today. `getScratch` takes one from the pool and hands it out with length 0; `putScratch` puts it back unless its capacity is over `maxScratch`, 65,536 passes.

> [!WARNING]
> A pool may hand out again what went back: a scratch grown for one long run would keep holding its memory, 6.8 MB at 65,536 passes.
