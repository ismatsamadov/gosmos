---
slug: a-station-type
title: A Station Type
type: code
position: 0
difficulty: 3
topics: structs, types
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
Nkem opens the network spreadsheet: Lantern Atoll appears twice, and Cape Wren's mask is wrong. "Let's model the network in code. Start with Port Ellery."

A struct type groups named fields into one value. A struct literal sets fields by name and leaves the rest at their zero values. `%+v` prints each field with its name:

```go
type Customer struct {
	Name       string
	Satellites int
}

c := Customer{Name: "Furrowlight", Satellites: 2}
fmt.Printf("%+v\n", c)
// {Name:Furrowlight Satellites:2}
```

## Your task

1. Define a `Station` type with the fields `Code` and `Name` (`string`), then `LatDeg` and `MaskDeg` (`float64`), in that order.
2. In `main`, replace the four loose variables with one `Station` value for Port Ellery.
3. Print it with `%+v`.
