---
slug: every-station-every-block
title: Every Station, Every Block
type: code
position: 6
difficulty: 5
topics: loops
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Ingrid needs to staff Pipit-3's passes, so the planner wants a grid: one row per station, one column per two-hour block of the day, and `#` where the station sees Pipit-3.

A loop inside another runs all its rounds for every round of the outer one. `fmt.Print` adds no newline, so a row stays on one line until `fmt.Println()` ends it:

```go
for dish := 1; dish <= 2; dish++ {
	fmt.Print("dish ", dish, " heaters:")
	for heater := 1; heater <= 3; heater++ {
		fmt.Print(" on")
	}
	fmt.Println()
}
```

## Your task

Inside the station loop in `main`, add a loop over the 12 blocks that prints `#` when `visible(station, block)` is true and `.` when it isn't.
