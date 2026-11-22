---
slug: the-round-trip
title: The Round Trip
type: code_tests
position: 8
difficulty: 5
topics: errors, structs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-22T10:19:20Z
---
Yusra's tools match the van's replies to her bookings by ID, so an accepted booking must come back with its ID exactly as she sent it. A booking without an ID could never be matched, so it must be refused before anything is written.

A round trip decodes, checks, then encodes: a booking that fails the check never reaches the encoder, and the caller gets a reason instead of a reply. The tags on `Booking` keep each key's name, and `omitempty` drops an empty note.

```go
var s struct{ Code string }
if err := json.Unmarshal(data, &s); err != nil {
    return err
}
if s.Code == "" {
    return errors.New("station has no code")
}
```

## Your task

Implement `confirm` with the read-only `Booking` in `booking.go`: decode `data`, refuse a booking whose ID is missing, empty or `null` with an error, and return the booking encoded as JSON. Return decoding errors too, and no bytes with any error.
