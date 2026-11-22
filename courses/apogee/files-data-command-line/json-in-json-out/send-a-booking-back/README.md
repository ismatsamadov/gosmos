---
slug: send-a-booking-back
title: Send a Booking Back
type: code_tests
position: 2
difficulty: 4
topics: errors, structs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
The van has checked bk-2041 against the calibration plan and must send it back to Yusra as JSON at the next signal. A reply pieced together with `fmt.Sprintf` would break on the first note that holds a quote.

`json.Marshal(v)` returns `v` as JSON bytes, with keys named by the tags and quotes and other special characters inside strings escaped. It returns an error for a value it cannot write, such as a time in the year 20270: a time goes out as text with a four-digit year.

```go
m := map[string]string{"dish": `LG "big"`}
data, err := json.Marshal(m)
fmt.Println(string(data), err)
// {"dish":"LG \"big\""} <nil>
```

## Your task

Implement `encodeBooking` with `json.Marshal`. Return the bytes on success, and nil with the error when `Marshal` fails: one test hands it a booking whose year was typed as 20270.
