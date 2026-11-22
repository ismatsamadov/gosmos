---
slug: write-the-station-clock
title: Write the Station Clock
type: code_tests
position: 1
difficulty: 6
topics: formatting, testing, time
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Paz heads every page of the calibration log with its time, and each laptop writes the time its own way: one prints 9/21/27 2:02 PM, another 21.09.2027 14:02. She wants one heading on every laptop.

`t.Format(layout)` writes a time the way a layout shows. A Go layout is the reference time, Monday 2 January 2006 at 15:04:05, written the way you want yours: each part of it stands for the same part of your time.

```go
t := time.Date(2027, time.September, 20,
	6, 5, 0, 0, time.UTC)
fmt.Println(t.Format("Monday 02/01, 3:04 PM"))
// Monday 20/09, 6:05 AM
```

## Your task

Implement `heading` to write `t` as these two headings show: short weekday and month, the day without a leading zero, a 24-hour clock to the minute, then `UTC`. gctl only hands it times in UTC.

```text
Tue 21 Sep 2027, 14:02 UTC
Fri 1 Oct 2027, 09:30 UTC
```
