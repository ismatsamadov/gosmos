---
slug: sorting-the-day
title: Sorting the Day
type: code_tests
position: 7
difficulty: 5
topics: functions, slices, sorting
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:28Z
---
Predictions arrive from passplan in whatever order its orbit math finishes them. The Board needs them by AOS.

`slices.SortFunc(s, compare)` sorts `s` in place. `compare(a, b)` returns a negative number when `a` goes first, a positive one when `b` does, and 0 for a tie. You can write it right in the call as a function literal, a function with no name:

```go
slices.SortFunc(stations, func(a, b Station) int {
	return cmp.Compare(a.MaskDeg, b.MaskDeg)
})
```

`cmp.Compare` compares two numbers or two strings that way, and `a.AOS.Compare(b.AOS)` compares two times.

## Your task

Make `sortByAOS` sort the passes by AOS. When two passes share an AOS, the lower station code goes first: `LAN` before `YST`. Add `"cmp"` and `"slices"` to the imports.

> [!NOTE]
> `SortFunc` may leave elements that compare as 0 in any order, so only a tie-break makes their order certain.
