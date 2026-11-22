---
slug: a-list-of-passes
title: A List of Passes
type: code_tests
position: 4
difficulty: 4
topics: slices, structs, testing
created: 2026-11-06T16:34:04Z
modified: 2026-11-22T10:19:20Z
---
Yusra's export for Wednesday 11 August is one JSON array holding every pass booked that day. Her tool writes a day without bookings as `[]`, an older tool as `null`, and the van's copy must say `[]` either way, because Yusra's checker reads `null` as a broken export.

`json.Unmarshal` fills a slice from a JSON array, one element per entry, and one bad entry makes the whole call return an error. JSON `null` leaves the slice nil, and `json.Marshal` writes a nil slice as `null` but an empty one as `[]`.

```go
var none []string
a, _ := json.Marshal(none)
b, _ := json.Marshal([]string{})
fmt.Println(string(a), string(b)) // null []
```

## Your task

Implement `decodeBookings`. Return the day's bookings in order; for `[]` or `null`, return an empty slice that is not nil. On an error, return nil with it, never the entries that did decode.
