---
slug: counting-by-satellite
title: Counting by Satellite
type: code_tests
position: 3
difficulty: 8
topics: maps, metrics, mutexes
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:38Z
---
November is settled by bookingd's records: 4 failed passes, all Brine-3's at Lantern Atoll. From now on Brinewatch's number is one cell of a table, since uplinkd counts every pass by satellite and by outcome, the metric's labels, with a counter for each pair of their values.

The counters live in a map that many goroutines read at once, while a new pair is rare. A `sync.RWMutex` lets readers share it under `RLock`; `Lock` waits for them to leave and lets one writer in. Between your `RUnlock` and your `Lock` another goroutine may add the same pair, so look again first:

```go
mu.RLock() // many readers at once
code := codes["Cape Wren"]
mu.RUnlock()

mu.Lock() // one writer, and no readers
codes["Lantern Atoll"] = "LAN"
mu.Unlock()
```

## Your task

`With` makes a new counter on every call, so every count is lost; `Snapshot` is already written. Keep each pair's counter in `v.m` under `[2]string{a, b}`: look it up under `RLock`; when it's missing, take `Lock`, look again, and only then add it, making the map at first use.
