---
slug: the-pass-that-moved
title: The Pass That Moved Itself
type: code
position: 5
difficulty: 5
topics: slices
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:19:30Z
---
Yusra starts her 09:10 message with a short summary, marked TL;DR: "Furrowlight says their 05:03 pass moved to 05:30 on the Board. Nobody moved it." You did, in a draft.

Slicing doesn't copy. `live[:2]` is a view onto the first two elements of the array `live` uses, so writing through the view changes `live` too:

```go
masks := []float64{5, 5, 3}
view := masks[1:]
view[0] = 9
fmt.Println(masks) // [5 9 3]
```

`slices.Clone(s)` returns a new slice with its own copy of the elements.

## Your task

Build `draft` with `slices.Clone`, so moving the pass in the draft leaves `live` alone. Add `"slices"` to the imports.
