---
slug: zero-values
title: Nothing Yet Is Still Something
type: code
position: 2
difficulty: 2
topics: types, variables
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
The scheduler creates a pass record before it knows anything about the pass. Until the record is filled in, it shows 0, false and empty text.

A variable declared without a value holds its type's zero value: `0` for numbers, `false` for a `bool` and `""`, empty text, for a `string`. To give an existing variable a new value, use `=`; `:=` is for declaring new ones.

```go
var gustKmh int
fmt.Println(gustKmh) // 0
gustKmh = 81
fmt.Println(gustKmh) // 81
```

## Your task

The scheduler has picked the station. Where the TODO is, assign `"PEL"` to `station`, and leave the other two variables as they are.

> [!NOTE]
> Maximum elevation, printed as `Max elevation`, is the highest a satellite climbs above the horizon during a pass, in degrees.
