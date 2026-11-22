---
slug: one-interface-for-bookings
title: One Interface for Bookings
type: code_tests
position: 0
difficulty: 7
topics: compiler-errors, interfaces
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T16:20:02Z
---
On Monday the team swapped the beta that was running for a new version, a deploy, and Furrowlight's twelve test bookings vanished with the old process's memory. A store on disk is coming, and the handlers mustn't care which store they get.

An interface type lists methods, and any type that has them can stand where the interface is expected:

```go
type Clock interface {
	Now() time.Time
}

type Poll struct {
	Clock Clock // a real clock, or a test's
}
```

## Your task

In `service.go`, add to `Store` the three methods of `MemStore` it lacks, then make the `Store` field a `Store`:

```go
Create(context.Context, Booking) (Booking, error)
List(context.Context) ([]Booking, error)
Delete(context.Context, string) error
```

Until then the tests don't build: Go refuses the test's failing store where a `*MemStore` is expected, in one long line:

```text
cannot use &failingStore{} (value of type
*failingStore) as *MemStore value in struct
literal
```
