---
slug: the-fleet-registry
title: The Fleet Registry
type: code
position: 0
difficulty: 3
topics: maps
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:19:30Z
---
Yusra forwards an angry email: Furrowlight's dashboard gave one of their passes to Brine-2, a satellite that never left the device meant to release it into orbit, its deployer. Nkem wants the fleet in a registry you look up by ID, not a list someone searches by hand.

A map stores values under keys. Write a literal with `key: value` pairs, add a key by assigning to it, and look a value up with its key in square brackets:

```go
masks := map[string]float64{
	"YST": 5,
	"SQT": 3,
}
masks["LAN"] = 8
fmt.Println(masks["SQT"]) // 3
```

## Your task

1. Add `thistle` to `operators`. Its operator, the organization responsible for it, is Kessick Polytechnic.
2. Replace each `"?"` with a lookup in `operators`, like the one in the first `Println`, so the program prints:

```text
pipit-1: Apogee
furrow-b: Furrowlight
thistle: Kessick Polytechnic
```

> [!WARNING]
> `var m map[string]string` declares a nil map: reading from it gives zero values, but adding a key panics. Start from a literal, even an empty one: `map[string]string{}`.
