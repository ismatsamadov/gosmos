---
slug: numbers-in-range
title: Numbers in Range
type: code_tests
position: 5
difficulty: 9
topics: reflection, time, types
created: 2026-11-12T03:14:42Z
modified: 2026-11-12T09:50:21Z
---
Friday 1 June. The form from Vorrel Island (VRL), an icebound island in the far south, says 0 dishes, and this morning a test deploy of passplan 2 read `Slots: 0` and planned nothing.

A tag's bound is text, and the field's type says how to read it. `CanInt` is true for the signed integer kinds, and `Int()` returns their value as an `int64`; `CanFloat` and `Float()` work the same for floats. A `time.Duration`'s kind is `int64`, so tell it by its type.

```go
v := reflect.ValueOf(90 * time.Second)
fmt.Println(v.Kind(), v.CanInt()) // int64 true
dur := reflect.TypeFor[time.Duration]()
fmt.Println(v.Type() == dur)      // true
```

## Your task

`checkRule` now cuts each rule at its `=` and hands `min` and `max` to `measure`, which fits nothing yet. Write it: for a duration, an int or a float, read `text` with `time.ParseDuration`, `strconv.ParseInt` or `strconv.ParseFloat`, and return `cmp.Compare`'s result, the value and the bound, each as its own type, and `true`. Any other type, or a bound that doesn't parse, returns `false`.

> [!NOTE]
> The problems print both values with `%v`, so a `time.Duration` prints as `15m0s`.
