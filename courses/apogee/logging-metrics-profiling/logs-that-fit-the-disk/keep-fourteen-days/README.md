---
slug: keep-fourteen-days
title: Keep Fourteen Days
type: code_tests
position: 5
difficulty: 8
topics: files, slices, sorting, time
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
Fourteen days of lines and 1 GiB per program at every station: that is Ingrid's budget for logs, and after each rotation the rotator keeps inside it.

First go the program's rotated files dated more than 14 days before today's UTC day; then the oldest go, by day and then by number, until the rest fit. The numbers must sort as numbers, not as text: `os.ReadDir` lists `.10.log` before `.9.log`. `slices.SortFunc` sorts by any rule, and `cmp.Or` returns the first comparison that isn't a tie:

```go
slices.SortFunc(ps, func(a, b PassView) int {
	return cmp.Or(a.AOS.Compare(b.AOS),
		strings.Compare(a.ID, b.ID))
})
```

## Your task

`prune` deletes nothing. Make it sort the files `rotatedFiles` lists, oldest first, then delete those dated more than `keepDays` days before `now`'s UTC day, then the oldest of the rest until they fit in `maxTotal` bytes.

> [!NOTE]
> `rotatedFiles` lists only the program's rotated files, so the current file, other programs' files and names that don't parse are never touched.
