---
slug: load-it-once
title: Load It Once
type: code_tests
position: 0
difficulty: 8
topics: goroutines, mutexes
created: 2026-11-11T12:51:42Z
modified: 2026-11-12T08:43:35Z
---
Each of the Board's hundred screens loads the network, the table of twenty stations and eighty satellites, from Port Ellery's file server. One load would serve them all, if the others waited for it.

`sync.Once` runs a function exactly once, however many goroutines call `once.Do` together: the first call runs it, and every other call waits until it has returned. Keep its results in fields beside the `Once`, and every caller reads the same ones, an error included.

```go
var once sync.Once
var gust int
read := func() int {
	once.Do(func() { gust = readMast("CWR") })
	return gust
}
```

## Your task

`Network` calls `loadNetwork` every time. Call it inside `t.once.Do` instead, keeping the network and the error in `t.net` and `t.err`, and return those two fields.

> [!NOTE]
> `loadNetwork` is the test's: a stand-in for the file server that counts its calls. A zero `Tables` is ready to use, since a zero `sync.Once` hasn't run yet.
