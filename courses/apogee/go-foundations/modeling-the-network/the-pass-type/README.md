---
slug: the-pass-type
title: The Pass Type
type: code_tests
position: 5
difficulty: 4
topics: structs, time
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
Nkem writes the next type on the whiteboard: a pass is a satellite, a station, an AOS and an LOS. "Everything else follows from those four."

A `time.Time` is one instant, date included, and a struct field can hold one like any other value. `later.Sub(earlier)` returns the `time.Duration` between two instants, even across midnight:

```go
sent := time.Date(2027, 4, 15,
	4, 15, 0, 0, time.UTC)
echoed := time.Date(2027, 4, 15,
	4, 15, 2, 0, time.UTC)
fmt.Println(echoed.Sub(sent)) // 2s
```

`time.Date` takes the year, month, day, hour, minute, second, nanosecond and time zone.

## Your task

Make `duration` return how long a pass lasts, from its AOS to its LOS.

> [!TIP]
> `b.Sub(a)` is `b` minus `a`. Put the later time first, or the duration comes out negative.
