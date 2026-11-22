---
slug: a-growing-schedule
title: A Growing Schedule
type: code
position: 2
difficulty: 3
topics: slices
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:07:33Z
---
Predictions arrive one at a time from passplan, the program that plans Apogee's passes, and the schedule should grow with each one. So far it stays empty.

`append` returns the slice with new elements added at the end. The variable you pass in keeps its old length, so assign the result back. A slice declared with `var` starts empty, and `append` works on it:

```go
var readings []float64
readings = append(readings, 7.91)
readings = append(readings, 7.90, 7.88)
fmt.Println(len(readings)) // 3
```

## Your task

In the loop, `_ =` throws away the slice that `append` returns; it's there only because Go won't build an `append` whose result goes unused. Keep the result in `passes` instead, so all three predictions land in the schedule.
