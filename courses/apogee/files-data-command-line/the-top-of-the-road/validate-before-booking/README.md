---
slug: validate-before-booking
title: Validate Before Booking
type: code_tests
position: 2
difficulty: 7
topics: errors, loops, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Hollis's sheet books Furrowlight's second satellite, Furrow-B, at station SQX, and no station has that code. One wrong line must stop the whole batch before gctl writes a plan from it.

Check every record before using any of them. A loop that returns at the first problem, with an error that says which record and why, lets the caller refuse the batch while nothing is written yet. A `time.Time` that was never set is the zero time, which `IsZero` reports:

```go
for i, r := range readings {
	if r.At.IsZero() {
		return fmt.Errorf("line %d: no time", i+1)
	}
}
```

## Your task

Implement `validateBookings`. Return an error for the first booking with no ID, no AOS (the moment its pass begins), a LOS (the moment it ends) that is not after its AOS, or a station that is not in `stations`. Name the booking by its ID, or by its number from 1 when it has none. Return nil for a valid batch.
