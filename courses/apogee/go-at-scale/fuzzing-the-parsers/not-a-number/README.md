---
slug: not-a-number
title: Not a Number
type: code_tests
position: 4
difficulty: 8
topics: conditionals, testing
created: 2026-11-11T12:51:44Z
modified: 2026-11-12T08:43:35Z
---
Saturday's other line, `batt_v=nan` from Thistle-2 at 03:36, sailed through the parser, and the pager's battery rule, `BattV < 6.6`, never fired.

`strconv.ParseFloat` reads "nan" as NaN, which is never below, above or equal to any number. The old check refuses a value when `v < lo || v > hi`, which is false for NaN, so NaN gets through. Write the check the other way round: `v >= lo && v <= hi` must be true for a value to pass. For NaN it is false, so NaN is refused, and so is an infinity:

```go
v, _ := strconv.ParseFloat("nan", 64)
fmt.Println(v < 6.6)                 // false
fmt.Println(!(v >= 6.6 && v <= 8.2)) // true
```

## Your task

`reading`, in `parse.go`, refuses a value if `v < lo || v > hi`. Refuse it whenever `v >= lo && v <= hi` is false instead. The battery and the temperature both go through `reading`.

> [!NOTE]
> `FuzzParseLine` runs your `checkLine` from the last lesson on Saturday's lines, which Joss forwarded and Mirela saved under `testdata`.
