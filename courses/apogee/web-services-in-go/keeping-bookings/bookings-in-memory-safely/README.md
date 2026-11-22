---
slug: bookings-in-memory-safely
title: Bookings in Memory, Safely
type: code_tests
position: 1
difficulty: 7
topics: maps, mutexes
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Until today a hidden placeholder kept the bookings. Nkem wants the service's own store, and asks one question: "On launch day two customers book one pass in the same second. Who gets it?"

A map is safe for many goroutines when every method holds one mutex while it uses it. The conflict check and the insert need the same lock, or two callers can both pass the check:

```go
func (t *Tally) Add(code string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts[code]++
	return t.counts[code]
}
```

## Your task

Write `MemStore`'s four methods. Each first returns `ctx.Err()` if the context is done, then locks `m.mu`.

1. `Create` gives the next ID, from `bk-4501`, unless a booking of the same satellite and station that starts at the same moment, the pass's AOS, is already kept: then it wraps `ErrConflict`.
2. `Get` and `Delete` wrap `ErrNotFound`; the comments beside the errors show both messages.
3. `List` returns a new slice, sorted by AOS, then ID.
