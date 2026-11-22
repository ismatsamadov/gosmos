---
slug: a-new-file-each-day
title: A New File Each Day
type: code_tests
position: 4
difficulty: 8
topics: files, logging, time
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
The night shift asks for "Tuesday's log", and Tuesday has to be one set of files, in UTC, whatever zone a clock's times come in.

So the rotator also rotates when the clock's UTC day is after the current file's: the old file is renamed under its own day, so `uplinkd-2028-12-04.1.log` holds the 4th's last lines, and on the 5th the count starts again at 1. An empty current file is never renamed; it only takes the new day. A day written as `2006-01-02` compares as text in date order:

```go
day := func(t time.Time) string {
	return t.UTC().Format(time.DateOnly)
}
z := time.FixedZone("UTC+8", 8*60*60)
a := time.Date(2028, 12, 5, 7, 59, 0, 0, z)
b := a.Add(time.Minute) // 00:00 UTC
fmt.Println(day(a), day(b), day(b) > day(a))
// 2028-12-04 2028-12-05 true
```

## Your task

Make `Write` rotate first when the clock's UTC day is after `r.day` and the file isn't empty; when it is empty, give it the new day and a count of 0 instead. A clock that steps back never rotates.
