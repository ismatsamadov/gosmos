---
slug: next-passes-first
title: Next Passes First
type: code_tests
position: 1
difficulty: 6
topics: slices, sorting, time
created: 2026-11-07T06:23:24Z
modified: 2026-11-07T06:23:24Z
---
The television fits twelve rows, and the night shift reads the top one first. But passplan, the service that predicts every pass, sends the booked passes in any order, finished ones too, and the page lists them as they come.

A template should only print what it is given, so prepare its data in Go. `slices.SortFunc` sorts with your comparison, and `cmp.Or` returns the first of its arguments that isn't zero, which gives a sort a second key:

```go
// Windiest first, then by station code.
slices.SortFunc(rs, func(a, b Weather) int {
	return cmp.Or(
		cmp.Compare(b.WindKmh, a.WindKmh),
		cmp.Compare(a.Station, b.Station))
})
```

## Your task

`nextPasses` in `board.go` returns `all` as it came, and `Refresh` already calls it. Make it return a new slice of the passes whose LOS, the moment the satellite sets, is after `now`: ordered by AOS, the moment it rises, then by booking ID, and at most 12.

> [!WARNING]
> Never sort `all` itself: that slice is passplan's.
