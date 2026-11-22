---
slug: ninety-nine-and-a-half
title: Ninety-Nine and a Half
type: code_tests
position: 0
difficulty: 8
topics: conversions, errors, metrics
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:38Z
---
Yusra starts the whiteboard with what customers ask first: how many of their passes were flown. She writes it as passes flown over passes booked, which ops calls availability, and under it Lantern Atoll's quarter: 1,177 flown of 1,199.

A share of two counts is a `float64`: convert both before dividing, or integer division rounds it down to 0. Floats don't refuse a zero count, either: 0/0 quietly becomes NaN, a value that spoils every sum it joins. So a count that can't be is refused with an error, where it is found:

```go
flown, booked := 0, 0
fmt.Println(float64(flown) / float64(booked))
// NaN
fmt.Println(7/2, float64(7)/2)
// 3 3.5
```

## Your task

Make `availability` return flown divided by booked, as a `float64`. Refuse counts that can't be with 0 and one of these errors, word for word:

```text
no booked passes
more flown than booked
```

> [!NOTE]
> Availability runs from 0 to 1: Lantern Atoll's 1,177 of 1,199 is 0.98165, or 98.17%.
