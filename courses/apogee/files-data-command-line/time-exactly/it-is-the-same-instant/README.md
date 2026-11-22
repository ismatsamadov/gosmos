---
slug: it-is-the-same-instant
title: It Is the Same Instant
type: code_tests
position: 2
difficulty: 6
topics: testing, time
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
The van's laptop logged this morning's test pass at `2027-09-22T07:14:00-07:00`, and the network's log has it at `2027-09-22T14:14:00Z`. Ingrid, on a crackling link: "Copy. One pass, two clocks. Log it in UTC."

An offset changes how an instant is written, not which instant it is. `t.UTC()` returns the same instant set to UTC, so it formats the way the network writes times. Watch the date as well as the hour: an evening in the van is already the next morning in UTC.

```go
t, err := time.Parse(time.RFC3339,
	"2027-09-22T20:30:00-07:00")
if err != nil {
	return err
}
fmt.Println(t.UTC().Format(time.DateOnly))
// 2027-09-23
```

## Your task

`networkTime` parses a timestamp and writes it back unchanged, offset and all. Make it return the same instant in RFC 3339 as the network logs it: in UTC, with a `Z`.
