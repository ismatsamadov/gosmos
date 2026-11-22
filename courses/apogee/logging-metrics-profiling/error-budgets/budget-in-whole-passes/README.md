---
slug: budget-in-whole-passes
title: The Budget in Whole Passes
type: code_tests
position: 1
difficulty: 8
topics: conversions, metrics
created: 2026-11-10T15:36:37Z
modified: 2026-11-10T15:36:37Z
---
Yusra, at her desk: "TL;DR: how many of Lantern Atoll's passes may we lose in a quarter and still keep our word?" Each station's promise, its availability target, is written in basis points, hundredths of a percent: 9950 means 99.50% of booked passes flown, and Lantern Atoll's is 9900.

What a target lets a station lose is its error budget, and a pass is never lost by halves. Integer arithmetic keeps it exact: multiply first, then divide, and the division rounds down to whole passes. A float can't hold 0.9995 exactly, and the rounding shows:

```go
t := 9995
f := 1 - float64(t)/10000
fmt.Println(f*10000, int(f*10000))
// 4.999999999999449 4
fmt.Println(10000 * (10000 - t) / 10000)
// 5
```

## Your task

`budget` works in floats. Make it work in integers: multiply `booked` by what the target leaves, `10000 - int(t)`, and only then divide by 10000.

> [!NOTE]
> Dividing first loses everything: `(10000 - 9950) / 10000` is already 0.
