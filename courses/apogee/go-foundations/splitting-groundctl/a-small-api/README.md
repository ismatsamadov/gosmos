---
slug: a-small-api
title: A Small API
type: code_tests
position: 3
difficulty: 6
topics: packages, structs
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
The telemetry parser has moved to `telemetry/`, and on Monday the downlink and pager teams start calling it. Nkem: "Show them what they need. Nothing more."

A package's API is its exported names. Exporting a struct type isn't the whole job: its fields follow the same rule, so other packages can use only the fields whose names start with a capital letter.

```go
// Booking is exported, as are SatID and Station.
type Booking struct {
	SatID   string
	Station string
	price   int // only this package sees it
}
```

## Your task

`main.go` and the tests are the other teams' code: they call `telemetry.Parse` and read a `Frame`'s `SatID`, `Seq` and `BattV`. In `telemetry/telemetry.go`, export `frame`, its fields and `parse`, doc comments included. Keep the helper `battery` unexported.
