---
slug: memory-in-use-now
title: Memory in Use Now
type: code_tests
position: 1
difficulty: 9
topics: profiling
created: 2026-11-10T15:36:39Z
modified: 2026-11-10T15:36:39Z
---
Monday afternoon. Ingrid puts a month of Lantern Atoll's `heap_bytes`, the memory its uplinkd reports using, on the Board: it climbs 4 MiB at every failed pass and falls back only when uplinkd restarts. Before anyone reads more into it, Nkem wants the number to mean memory in use.

`runtime.ReadMemStats` fills a `runtime.MemStats` with the program's memory numbers. `HeapAlloc` is the bytes of heap blocks not yet freed, garbage included until a collection frees it; `TotalAlloc` adds up the bytes of every block ever allocated, so it never falls. Right after `runtime.GC()`, `HeapAlloc` is what is still in use:

```go
var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Println(m.NumGC) // collections so far
// 0
```

## Your task

`inUse` returns `TotalAlloc`. Make it run `runtime.GC()` first, then return `HeapAlloc`.

> [!NOTE]
> `runtime.GC` holds its caller until a whole collection is done, so uplinkd collects only when someone asks for the number, as the pager does every five minutes.
