---
slug: many-readers-one-writer
title: Many Readers, One Writer
type: code_tests
position: 1
difficulty: 7
topics: maps, mutexes
created: 2026-11-01T07:36:47Z
modified: 2026-11-04T17:07:34Z
---
Every service reads the station-status cache, a stored copy kept for quick access, and only the poller writes to it. Yesterday someone changed the map `Snapshot` handed out, and the Board showed Cape Wren stowed while its dish was tracking.

A `sync.RWMutex` locks two ways: `RLock` for readers, which many goroutines can hold at once, and `Lock` for a writer, who waits to have the data alone. Returning a map shares it rather than copying it, so hand out a copy made with `maps.Clone`.

```go
func (b *Board) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.rows)
}
```

## Your task

Make `Get` take the read lock instead of the full one. Make `Snapshot` return a copy of `c.status`, taken under the read lock.
