---
slug: numbers-without-rounding
title: Numbers Without Rounding
type: code_tests
position: 6
difficulty: 5
topics: conversions, maps, testing, types
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Paz's quick viewer decodes any report into a map, to show kinds gctl has no type for yet. This morning it listed reports 9007199254741003 and 9007199254741005 under one ID that belongs to neither.

Decoded into `any`, a JSON number becomes a `float64`, which holds every whole number exactly up to 2^53^ (9,007,199,254,740,992) but rounds some beyond it. After `UseNumber`, a decoder gives `json.Number` values instead, the digits kept as text; their `Int64` method converts them exactly or returns an error.

```go
const big = "9007199254740993"
var plain, exact any
_ = json.Unmarshal([]byte(big), &plain)
dec := json.NewDecoder(strings.NewReader(big))
dec.UseNumber()
_ = dec.Decode(&exact)
fmt.Println(plain) // 9.007199254740992e+15
fmt.Println(exact) // 9007199254740993
```

## Your task

Implement `readReportID`: decode the report from `r` into a `map[string]any` with `UseNumber`, and return its `id` as an `int64`. Return an error when the `id` is missing or not a number, and `Int64`'s error for a fraction or an overflow.
