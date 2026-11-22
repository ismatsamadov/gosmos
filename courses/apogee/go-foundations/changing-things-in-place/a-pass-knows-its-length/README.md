---
slug: a-pass-knows-its-length
title: A Pass Knows Its Length
type: code
position: 3
difficulty: 4
topics: methods, time
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
The Board's code computes `p.LOS.Sub(p.AOS)` in four places. Nkem's review asks for `p.Duration()` instead: "Works. Now make it say what it means."

A method is a function with a receiver, written before its name. Inside, the receiver is the value the method was called on:

```go
func (o Orbit) RadiusKm() float64 {
	return 6371 + o.AltKm
}

pipit1 := Orbit{AltKm: 525, IncDeg: 97.5}
fmt.Println(pipit1.RadiusKm()) // 6896
```

## Your task

`main` already calls `p.Duration()`, so the program doesn't build yet: `p.Duration undefined (type Pass has no field or method Duration)`. Where the `TODO` is, add the method to `Pass`, returning the time from AOS to LOS as a `time.Duration`.
