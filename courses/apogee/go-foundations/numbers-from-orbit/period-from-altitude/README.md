---
slug: period-from-altitude
title: The Period of an Orbit
type: code_tests
position: 8
difficulty: 4
topics: constants, packages, types
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
Tomás opens his notebook to a single line. "Altitude in, period out," he says, "and all it takes is one constant, μ: Earth's mass times the gravitational constant."

The `math` package holds what formulas need, such as `math.Pi` and `math.Sqrt(x)`, the square root of `x`. Go has no operator for powers, so write a cube as `a * a * a`.

```go
// a Ystmark dish's radius, in meters
r := 3.7 / 2
fmt.Println(math.Pi * r * r) // its area: 10.75...
fmt.Println(math.Sqrt(6.25)) // 2.5
```

## Your task

Add `"math"` to the imports. Then make `periodMinutes` return the period, the time one orbit takes: 2π·√(a³/μ), which is in seconds, divided by 60. Here a is `earthRadiusKm + altKm`, the orbit's radius, and μ is `muKm3S2`.
