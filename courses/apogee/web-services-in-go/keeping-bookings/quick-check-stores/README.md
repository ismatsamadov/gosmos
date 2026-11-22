---
slug: quick-check-stores
title: Quick Check: Stores
type: quiz
position: 2
difficulty: 7
topics: interfaces, mutexes, slices
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Mirela stops at the backend desks, reads `store.go` over Nkem's shoulder, and has three questions, precise as ever, about your new store:

```go
type MemStore struct {
	mu       sync.Mutex
	bookings map[string]Booking
	next     int
}

func (m *MemStore) List(ctx context.Context) (
	[]Booking, error)
```

## Your task

Answer Mirela's three questions: which values can be a `Store`, why the conflict check and the insert share one lock, and why `List` hands out a new slice.
