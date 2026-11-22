---
slug: lost-updates
title: Lost Updates
type: code_tests
position: 0
difficulty: 7
topics: defer, mutexes
created: 2026-11-01T07:36:47Z
modified: 2026-11-06T16:51:27Z
---
Your first on-call week opens with a page: the frame counter says Pipit-2 sent 1,861 frames this hour, and the station logs say 2,000.

When goroutines change the same map at once, updates get lost, or the program crashes. A `sync.Mutex`, ready to use as its zero value, lets one goroutine in at a time: `Lock` waits until nobody else holds it, and `Unlock` lets the next one in. Calling `defer mu.Unlock()` right after `Lock` releases it however the function returns.

```go
// Post adds a row; b.mu guards b.rows.
func (b *Board) Post(row string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.rows = append(b.rows, row)
}
```

## Your task

`Add` does nothing yet. Make it lock `c.mu`, unlock it with `defer`, and count one frame for `sat`.
