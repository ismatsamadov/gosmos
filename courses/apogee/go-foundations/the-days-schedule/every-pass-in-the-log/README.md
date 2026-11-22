---
slug: every-pass-in-the-log
title: Every Pass in the Log
type: code
position: 1
difficulty: 3
topics: loops, slices
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:19:30Z
---
Ingrid's pass log numbers each pass from 1. Today's log has three passes; Friday's will have forty. Each pass ID names its satellite, station and orbit number, separated by slashes.

`for i, v := range s` runs once per element, in order, with the index in `i` and a copy of the element in `v`. Indexes start at 0:

```go
bands := []string{"UHF", "S", "X"}
for i, b := range bands {
	fmt.Printf("band %d is %s\n", i, b)
}
// band 0 is UHF, band 1 is S, band 2 is X
```

## Your task

Replace the single `Println` with a `range` loop over `ids` that prints every pass, numbered from 1, as `1 pipit-1/PEL/04182`.
