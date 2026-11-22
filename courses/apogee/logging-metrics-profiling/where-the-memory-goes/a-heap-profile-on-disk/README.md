---
slug: a-heap-profile-on-disk
title: A Heap Profile on Disk
type: code_tests
position: 0
difficulty: 9
topics: files, profiling
created: 2026-11-10T15:36:39Z
modified: 2026-11-11T17:33:38Z
---
Lantern Atoll's runbook, the steps the night shift follows when something breaks, still ends with "restart uplinkd". Nkem gives it a new first step: save uplinkd's memory to a file.

A program keeps its data in the heap, the memory it asks for as it runs, and Go's garbage collector frees whatever nothing can reach any more. A heap profile samples the blocks still in use and records the line that asked for each. `pprof.Lookup("heap")` returns it, and its `WriteTo` writes it as of the last collection, so run `runtime.GC()` first. Every profile writes the same way; this one writes each goroutine's stack as text:

```go
p := pprof.Lookup("goroutine")
p.WriteTo(os.Stdout, 1) // 1: text for people
// goroutine profile: total 1
```

## Your task

`heapProfile` only creates the file. Make it run `runtime.GC()`, write the heap profile into the file with `WriteTo(f, 0)`, and close the file, returning the first error.

> [!TIP]
> `go tool pprof -top -inuse_space heap.prof`, run in a terminal in the lesson folder, lists what holds the memory. uplinkd's debug port serves its heap profile too.
