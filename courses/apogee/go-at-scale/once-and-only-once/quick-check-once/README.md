---
slug: quick-check-once
title: Quick Check: Once
type: quiz
position: 3
difficulty: 8
topics: functions, goroutines
created: 2026-11-11T12:51:43Z
modified: 2026-11-11T12:51:43Z
---
Tuesday evening, before the Board restarts on the new table, Nkem has three questions about `sync.Once` and its helpers.

The first is about this loader, which two screens, A and B, call at the same moment:

```go
var once sync.Once
var net *Network
get := func() *Network {
	once.Do(func() { net = slowLoad() })
	return net
}
```

## Your task

Answer Nkem's questions: what B gets while A's call is loading, what a `sync.OnceValues` function returns after a failed load, and what happens when the function given to `once.Do` calls it again.
