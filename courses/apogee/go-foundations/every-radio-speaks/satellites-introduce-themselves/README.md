---
slug: satellites-introduce-themselves
title: Satellites That Introduce Themselves
type: code
position: 3
difficulty: 4
topics: formatting, interfaces, methods
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:27Z
---
The pass logs print satellites as `{pipit-1 Pipit-1 {525 97.5}}`, and Ingrid wants `Pipit-1 (525 km)`.

Last week a `String` method made `Severity` print as a word. Here's why that worked: `fmt` declares an interface, `Stringer`, whose one method is `String() string`. Before printing a value, `fmt` checks whether it satisfies `Stringer`, and if it does, prints what `String` returns:

```go
type Pass struct{ SatID, Station string }

func (p Pass) String() string {
	return p.SatID + "@" + p.Station
}

fmt.Println(Pass{"thistle", "LAN"}) // thistle@LAN
```

## Your task

Give `Satellite` a `String` method so `main` prints each satellite as its name and altitude, like `Pipit-1 (525 km)`. Leave `main` as it is.
