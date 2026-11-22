---
slug: tonights-passes
title: Tonight's Passes
type: code_tests
position: 0
difficulty: 5
topics: slices, sorting, time
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T17:14:44Z
---
Ingrid's first question every night is what flies over, and when. `data.go` now holds Friday night's passes, listed station by station, as the schedule file has them.

`slices.Clone` copies a slice, so sorting the copy leaves the original alone, and `slices.SortFunc` sorts with a comparison you supply:

```go
first := slices.Clone(times)
slices.SortFunc(first, func(a, b time.Time) int {
	return a.Compare(b)
})
```

## Your task

1. Give `Shell` a `passes []Pass` field, set to `tonight()` in `newShell`.
2. Write `cmdPasses`. It prints one line per pass, in time order, and leaves `s.passes` as it was. A line holds the moments the satellite rises into view and sets, which ops calls AOS and LOS, then the station, the satellite and the highest it climbs, in degrees:

```text
04:12-04:21 PEL pipit-1 max 61.4
```

3. Add `passes` to the table, so `help` lists it as:

```text
passes    list tonight's passes
```
