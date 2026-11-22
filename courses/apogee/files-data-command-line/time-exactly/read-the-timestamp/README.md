---
slug: read-the-timestamp
title: Read the Timestamp
type: code_tests
position: 0
difficulty: 6
topics: errors, testing, time
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Yusra's booking export reached the van at the last stop with a signal. Each booking gives the moment its satellite rises into view, called AOS, as text such as `2027-09-21T14:02:10Z`, and gctl can't compare passes until that text is a time.

`time.Parse(layout, text)` reads text written in a layout and returns a `time.Time` and an error. The layout `time.RFC3339` matches these timestamps: the date, `T`, the time, then `Z` for UTC or the hours and minutes the clock runs behind or ahead of UTC, its offset, such as `-07:00`.

```go
d, err := time.Parse(time.DateOnly, "2027-09-20")
if err != nil {
	return err
}
fmt.Println(d.Weekday()) // Monday
```

## Your task

Implement `parseStart` with `time.Parse` and `time.RFC3339`. Return the time, or the zero time and the error: a start without its `Z` or offset is not exact.
