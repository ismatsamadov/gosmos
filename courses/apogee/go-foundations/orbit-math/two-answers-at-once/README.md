---
slug: two-answers-at-once
title: Two Answers at Once
type: code_tests
position: 2
difficulty: 3
topics: functions, multiple-returns
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
The scheduler stores times as minutes after midnight: Pipit-1's AOS over Port Ellery is `252`. The Board wants hours and minutes.

A function can return several values. List their types in parentheses, return them separated by commas, and give each one a name where you call the function:

```go
// window returns a pass's AOS and LOS,
// in minutes after midnight.
func window(aosMin, lengthMin int) (int, int) {
	return aosMin, aosMin + lengthMin
}

aos, los := window(252, 9) // 252 and 261
```

## Your task

Make `hoursMinutes` return the hours and the minutes in `total`: `252` is 4 hours and 12 minutes.
