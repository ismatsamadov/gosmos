---
slug: every-pass-of-the-quarter
title: Every Pass of the Quarter
type: code_tests
position: 7
difficulty: 9
topics: iterators, slices
created: 2026-11-11T21:35:30Z
modified: 2026-11-11T21:35:30Z
---
Friday 30 March. Furrowlight books its summer campaign today, and Hollis's program wants every pass of Furrowlight's twenty-six satellites in the third quarter, read from bookingd a pass at a time.

Filter before the costly step, not after it. A loop over `Catalog` that keeps only Furrowlight's passes still predicts all eighty satellites, fifty-four of them for nothing. Choose the satellites first, then predict only those, a day at a time:

```go
var south []Site
for _, s := range n.Sites {
	if s.LatDeg < 0 {
		south = append(south, s)
	}
}
```

`Catalog` makes the days of any `Network` it's given, and a network's `Elements` list its `Satellites` in the same order.

## Your task

`Bookable` yields nothing. Make it yield every pass in `w` of the satellites whose `Operator` is `operator`, at the stations that can serve them, in `byAOS` order and a day at a time, without asking `predict` about any other satellite.
