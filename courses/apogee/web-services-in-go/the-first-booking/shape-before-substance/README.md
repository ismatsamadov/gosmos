---
slug: shape-before-substance
title: Shape Before Substance
type: code_tests
position: 6
difficulty: 7
topics: errors, runes, switch
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Hollis's program sends bookings without `los`, the moment the satellite sets out of the station's view, which ops calls LOS. The service keeps them anyway, so nobody knows when those passes end.

A body that decodes can still make no sense. Check its shape field by field before anything else reads it, and answer `422 Unprocessable Entity` naming the first rule it breaks, so the caller fixes one thing at a time. Count text in characters, not bytes: `len` counts bytes, and an `Ñ` takes two.

```go
note := "AÑIL"
fmt.Println(len(note))                    // 5
fmt.Println(utf8.RuneCountInString(note)) // 4
```

## Your task

`checkShape` returns nil. Return `invalidBooking` with the message of the first rule `b` breaks, in this order, or nil when it breaks none:

```text
id is chosen by the service
sat_id is required
station is required
aos and los are required
los must be after aos
note is longer than 200 characters
```
