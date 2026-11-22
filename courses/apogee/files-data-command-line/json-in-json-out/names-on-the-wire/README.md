---
slug: names-on-the-wire
title: Names on the Wire
type: code_tests
position: 1
difficulty: 4
topics: structs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
At the next signal Yusra checks the van's copy of bk-2042 and finds no satellite in it. Her export writes `sat_id`, Go's field is `SatID`, and ignoring case does not make an underscore match.

A struct tag, a string in backquotes after a field's type, can name the field's JSON key: `json:"sat_id"`. Decoding reads that key into the field, and encoding writes it. Without a tag, encoding writes the Go name, so `SatID` would go out as `"SatID"`.

```go
type Reading struct {
    Dish    string  `json:"dish"`
    ElevDeg float64 `json:"elev_deg"`
}
```

## Your task

Tag every field of `Booking` in `booking.go` with its key in the export: `id`, `sat_id`, `station`, `aos` and `los`. The tests decode one of Yusra's bookings and encode it back, and both directions must use exactly her keys.

> [!WARNING]
> Write no space after the colon. With a space there, the tag is not valid, and the json package ignores it.
