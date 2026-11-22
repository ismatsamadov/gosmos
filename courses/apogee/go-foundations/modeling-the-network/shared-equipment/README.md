---
slug: shared-equipment
title: Shared Equipment
type: code
position: 6
difficulty: 5
topics: compiler-errors, structs
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
The mobile test van carries a 1.2 m dish to sites to check and correct the dishes' pointing, a job called calibration. Its record should read like a station's. Nkem: "Same `Dish` fields, written once."

Embed a struct by writing its type with no field name. Its fields are promoted: you use them as if the outer struct declared them, even as keys in a literal (new in Go 1.27). A literal can also set the embedded struct whole, by its type's name:

```go
type Radio struct{ FreqMHz float64 }
type Beacon struct {
	SatID string
	Radio // embedded
}

p := Beacon{Radio: Radio{FreqMHz: 437.25}}
t := Beacon{SatID: "thistle", FreqMHz: 435.9}
fmt.Println(p.FreqMHz, t.FreqMHz) // 437.25 435.9
```

## Your task

The program doesn't build: `report.go`, which you can't edit, reads `yst.DiameterM`. Embed `Dish` in `Station` and in `Van`, in place of their `Equipment` fields, and update the two literals to match.
