---
slug: structs-as-parameters
title: Structs as Parameters
type: code_tests
position: 1
difficulty: 3
topics: formatting, functions, structs
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
The Board will list each station on one tidy line. Ingrid has one request: "No minus signs. Nobody in the Fishbowl reads them at a glance. Write N or S."

A function can take a whole struct as one parameter and read its fields with a dot:

```go
// highMask reports whether a station's mask
// is above 6 degrees.
func highMask(s Station) bool {
	return s.MaskDeg > 6
}
```

## Your task

Make `stationLine` return a station's code, name, latitude and mask, as `PEL Port Ellery (46.9 N, mask 5.0)`. `LatDeg` is negative in the south, but the line shows its size and an `S`: Salar Quieto is `23.4 S`.
