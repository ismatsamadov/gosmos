---
slug: sort-once-at-the-end
title: Sort Once at the End
type: code_tests
position: 6
difficulty: 9
topics: profiling, sorting
created: 2026-11-10T15:36:39Z
modified: 2026-11-10T15:36:39Z
---
Friday morning, the profile's last fix: passplan sorted its whole plan again after adding each pass, so 600 passes meant 600 sorts. Nkem: "Same order, one sort. What does sorting the first 599 again tell us?"

Sorting after every insert sorts a nearly sorted list over and over; appending everything first and sorting once does that work one time. `slices.SortFunc` sorts with a function that compares two elements, answering negative, zero or positive, and a tie on one key can fall through to the next:

```go
ids := []string{"pipit-1", "brine-3", "furrow-a"}
slices.SortFunc(ids, func(a, b string) int {
	return cmp.Or(cmp.Compare(len(a), len(b)),
		strings.Compare(a, b))
})
fmt.Println(ids)
// [brine-3 pipit-1 furrow-a]
```

## Your task

Make `buildDay` append every group's passes to one new slice, then sort it once with `slices.SortFunc` and `earlier`, in `passplan.go`, which orders passes by the second they start, their AOS, then by ID.

> [!NOTE]
> The groups belong to the caller: the day is a slice of its own, so sorting it leaves them as they were.
