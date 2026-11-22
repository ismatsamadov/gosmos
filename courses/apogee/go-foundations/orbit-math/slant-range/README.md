---
slug: slant-range
title: Slant Range
type: code_tests
position: 5
difficulty: 4
topics: functions
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
Before a pass, Tomás checks how far the signal travels: the slant range, the straight-line distance from dish to satellite. For Pipit-1 it is 525 km overhead and over 2,600 km at the horizon.

The `math` package has what you need: `math.Sqrt`, and `math.Sin` and `math.Cos`, which take radians, an angle measure where π radians is 180 degrees. Go has no power operator, so square a value by multiplying it by itself:

```go
// diagonalKm returns the distance across
// two sides at a right angle.
func diagonalKm(aKm, bKm float64) float64 {
	return math.Sqrt(aKm*aKm + bKm*bKm)
}
```

## Your task

Make `slantRangeKm` return √((R+h)² − (R·cos e)²) − R·sin e, where R is `earthRadiusKm`, h is `altKm`, and e is `elevDeg` converted with `radians`.
