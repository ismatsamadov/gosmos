---
slug: two-uploads-per-station
title: Two Uploads per Station
type: code_tests
position: 6
difficulty: 9
topics: context, goroutines, maps, mutexes
created: 2026-11-11T12:51:42Z
modified: 2026-11-11T12:51:42Z
---
Thursday 15 March. Calidra Atoll (CAL), a new station on an atoll near the equator, goes live on 2 April, and this morning three uploads to its dish loader ran at once over the host's thin line and timed out together.

A host's line carries two uploads at a time, so each station needs slots of its own: a map from station code to `Slots`, each made the first time its station is asked for. Goroutines share the map, so a `sync.Mutex` guards it. Hold the lock only to find or make the slots, never while an upload runs.

```go
mu.Lock()
log, ok := logs[code]
if !ok {
	log = newLog(code) // first use
	logs[code] = log
}
mu.Unlock()
log.Write(line) // outside the lock
```

## Your task

`Do` runs every upload at once. Find or make the station's `Slots`, with `u.per` slots, under `u.mu`, then hold one while `upload` runs: return `Acquire`'s error if the context ends first, or what `upload` returns.

> [!WARNING]
> One `Slots` shared by every station would make Varsk Island wait for Calidra Atoll.
