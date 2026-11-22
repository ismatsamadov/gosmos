---
slug: in-time-order
title: In Time Order
type: code_tests
position: 6
difficulty: 7
topics: logging, sorting, time
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Friday morning. The four programs' lines come from four computers, each log in its own order, and the night shift wants one timeline: bookingd's booking on Monday, the station's failure on Wednesday, everything between in its place.

Join the logs, then sort by time. `slices.SortStableFunc` sorts with a comparison function and keeps elements that compare equal in the order they came: headquarters' `lines delivered` and Lantern Atoll's `lines queued` share 01:48:30, and the one from the log given first stays first. `Time.Compare` returns the -1, 0 or +1 such a function needs:

```go
s := []string{"bb", "a", "cc", "d"}
slices.SortStableFunc(s, func(x, y string) int {
	return len(x) - len(y)
})
fmt.Println(s) // [a d bb cc]
```

## Your task

`merge` joins the logs into one new slice. Sort that slice by each line's `Time` with `slices.SortStableFunc` before returning it.

> [!WARNING]
> `slices.SortFunc` gives the same order for different times, but it may swap lines whose times are equal.
