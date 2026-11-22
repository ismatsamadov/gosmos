---
slug: nobody-home
title: Nobody Home
type: code_tests
position: 2
difficulty: 5
topics: pointers, slices
created: 2026-11-01T07:36:46Z
modified: 2026-11-01T07:36:46Z
---
Someone typed `XYZ` into the stow tool, and it reported a station with no name as stowed. A lookup that finds nothing has to say so.

A function can return a pointer, or `nil` when there is nothing to point to, and the caller checks for `nil` before using it. `&sats[i]` points at the element inside the slice:

```go
func newest(sats []Satellite) *Satellite {
	if len(sats) == 0 {
		return nil
	}
	return &sats[len(sats)-1]
}
```

## Your task

Complete `find` so it returns a pointer to the station in `stations` whose `Code` is `code`, or `nil` if there is none. A change made through the pointer must show in the slice.

> [!WARNING]
> In `for _, s := range stations`, `s` is a copy, so `&s` points at the copy, not into the slice.
