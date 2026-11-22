---
slug: the-background-refresher
title: The Background Refresher
type: code_tests
position: 3
difficulty: 8
topics: goroutines, mutexes
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
`main.go` now starts a refresher, a goroutine that loads the passes again every minute while the night shift types. Nkem asks: "What does the cache hold if `refresh` and `passes` run at the same moment?"

A `sync.Mutex` lets one goroutine through at a time: `Lock` waits until it is free, and a deferred `Unlock` frees it when the method returns. A goroutine that holds it and locks it again waits forever.

```go
func (b *board) add(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = append(b.lines, line)
}
```

## Your task

1. Give `passCache` a `mu sync.Mutex` field.
2. Move the body of `refresh` into `reload`, which expects `c.mu` held, and make `refresh` lock and call it.
3. Make `passes` lock, call `reload` when it must, and return a copy of `c.cached`.
