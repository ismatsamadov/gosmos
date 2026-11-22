---
slug: order-by-several-keys
title: Order by Several Keys
type: code_tests
position: 0
difficulty: 8
topics: functions, sorting
created: 2026-11-11T12:51:43Z
modified: 2026-11-12T08:43:35Z
---
Yesterday at 10:29:34, Brine-3 and Pipit-3 rose over Port Ellery in the same second, with one dish between them. Nkem: "What decides which one gets it?"

A comparison returns a negative number, zero or a positive one. Zero means equal, and a sort may put equal values in either order. `cmp.Or` returns the first of its arguments that isn't zero, so a second key settles only a tie:

```go
func byMask(a, b Station) int {
	return cmp.Or(
		cmp.Compare(a.MaskDeg, b.MaskDeg),
		strings.Compare(a.Code, b.Code),
	)
}
```

## Your task

`byAOS` compares only the AOS, the moment a pass rises into view, so those two passes compare equal. Make the pass ID the second key, with `strings.Compare`.

> [!WARNING]
> The AOS stays the first key: an earlier pass comes first, whatever its ID.
