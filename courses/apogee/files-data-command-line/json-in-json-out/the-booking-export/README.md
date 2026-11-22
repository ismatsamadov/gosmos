---
slug: the-booking-export
title: The Booking Export
type: code_tests
position: 0
difficulty: 4
topics: errors, structs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
At 1,900 m, a few minutes of signal bring one download: Yusra, who looks after Apogee's customers, has exported every pass they reserved at Salar Quieto, each one a booking. Calibration must never take a booked pass, so each booking has to reach the van exactly as she sent it.

`json.Unmarshal(data, &v)` fills `v` from JSON text. A key fills the exported field whose name matches it, ignoring case, and a key with no such field is skipped. Broken JSON gives a `*json.SyntaxError`, a value of the wrong kind a `*json.UnmarshalTypeError`. Times written as `2027-08-11T14:02:10Z` fill `time.Time` fields.

```go
var s struct{ Code string }
err := json.Unmarshal([]byte(`{"code":"LG"}`), &s)
fmt.Println(s.Code, err) // LG <nil>
```

## Your task

Implement `decodeBooking` for the read-only `Booking` in `booking.go`. On any error, return the zero `Booking` with it: a half-filled booking must not reach the van's plan. The export's `sat_id` matches no field yet; that is the next lesson.
