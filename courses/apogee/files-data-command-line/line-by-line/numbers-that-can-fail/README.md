---
slug: numbers-that-can-fail
title: Numbers That Can Fail
type: code_tests
position: 5
difficulty: 4
topics: conversions, errors, testing
created: 2026-11-06T16:34:04Z
modified: 2026-11-06T17:07:41Z
---
One reading from 6 July was typed in by hand, and its elevation says `41..70`: a decimal point typed twice. A careless parser turns it into 0°, La Grande pointing flat at the horizon, and Paz would chase an error that never happened.

`strconv.ParseFloat(text, 64)` turns text such as `41.70` or `-0.35` into a `float64`. When the text is not a number, it returns 0, a placeholder rather than a reading, and an error that names the text:

```go
volts, err := strconv.ParseFloat("7.9l", 64)
fmt.Println(volts) // 0
fmt.Println(err)   // prints:
```

```
strconv.ParseFloat: parsing "7.9l": invalid syntax
```

## Your task

Implement `parseElevation`. Return the number for valid text, negative values included: near the horizon a reading can sit a few tenths of a degree below zero. For anything else, return 0 and the error from `ParseFloat`.
