---
slug: below-zero-in-the-radome
title: Below Zero in the Radome
type: code_tests
position: 1
difficulty: 7
topics: bytes, conversions
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
dishd's screen says the air inside Dish 1's radome is at 6,541.1 °C. Sunniva, standing in it with a flashlight, can see each breath: it is −12.5 °C.

The temperature is a signed 16-bit field, in tenths of a degree, from −3276.8 to 3276.7. Read as a `uint16`, its two bytes are never negative. Converting to `int16` reads the same 16 bits with a sign; then convert to `float64` and divide, in that order:

```go
var u uint16 = 0xfffe
fmt.Println(u, int16(u)) // 65534 -2
fmt.Println(float64(int16(u)) / 4)
// -0.5
```

## Your task

`tempC` reads bytes 8 and 9 of a status body as an unsigned number. Convert them through `int16` before converting to `float64` and dividing by 10, so `ff 83` gives −12.5.

> [!WARNING]
> Convert to `float64` before dividing: `int16(u) / 10` is integer division, and −125 / 10 is −12.
