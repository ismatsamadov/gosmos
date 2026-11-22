---
slug: the-van-s-fixed-offset
title: The Van's Fixed Offset
type: code_tests
position: 4
difficulty: 6
topics: testing, time
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Paz wants the calibration log to show the van's own clock beside UTC. The salt flat's clock is simple: seven hours behind UTC all year, with no daylight saving.

A named time zone needs the Mac's time zone files, which change when a region's laws do. `time.FixedZone(name, secondsEast)` needs none: it makes a location that is always the same offset from UTC, counted in seconds east, so west is negative. `t.In(loc)` returns the same instant shown in that location.

```go
plusTwo := time.FixedZone("", 2*60*60)
noon := time.Date(2027, time.September, 23,
	12, 0, 0, 0, plusTwo)
fmt.Println(noon.UTC().Format("15:04")) // 10:00
```

## Your task

Implement `vanTime` to return `t` shown at a fixed offset of `secondsEast` seconds: the same instant, on that clock. gctl calls it with `-7*60*60` for the van.

> [!WARNING]
> `t.Add` moves the instant itself. To show one moment on another clock, change only its location.
