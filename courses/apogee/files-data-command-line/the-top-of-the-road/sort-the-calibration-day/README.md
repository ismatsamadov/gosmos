---
slug: sort-the-calibration-day
title: Sort the Calibration Day
type: code_tests
position: 3
difficulty: 7
topics: slices, sorting, testing, time
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Two of Thursday's passes start at 15:40:05 UTC, one on each dish, and Paz wants them in the same order every time she prints the plan. The rest of gctl still reads the bookings in the order they arrived.

`slices.SortFunc` sorts a slice in place by a function that returns a negative number, zero or a positive one. `cmp.Or` returns its first argument that isn't zero, so a second key settles ties. Sort a copy made by `slices.Clone` to leave the caller's slice as it was. `time.Time`'s `Compare` compares instants, whatever offset each time was written with:

```go
sorted := slices.Clone(readings)
slices.SortFunc(sorted, func(a, b Reading) int {
	return cmp.Or(
		cmp.Compare(a.ElevDeg, b.ElevDeg),
		strings.Compare(a.Dish, b.Dish),
	)
})
```

## Your task

Implement `orderedPasses`: return a copy of `passes` sorted by AOS, the moment each pass begins, and by booking ID when two begin together. Leave `passes` unchanged.
