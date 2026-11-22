---
slug: double-booked
title: Double-Booked at Ystmark
type: code_tests
position: 9
difficulty: 6
topics: loops, slices, time
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Sunniva, 18:02: "Dish 2 radome repair until Friday. One dish tonight. Check bookings." Every pair of passes needs checking, and each pair only once.

To visit every pair once, start the inner loop one past the outer index:

```go
codes := []string{"YST", "PEL", "SQT"}
for i := range codes {
	for j := i + 1; j < len(codes); j++ {
		fmt.Println(codes[i], codes[j])
	}
} // YST PEL, YST SQT, PEL SQT
```

Two passes overlap when each starts before the other ends; `t.Before(u)` reports whether time `t` is earlier than `u`. Passes that only touch, one's LOS at the other's AOS, don't overlap.

## Your task

Make `overlaps` return one line per pair of passes at the same station whose times overlap, as `furrow-b@YST overlaps pipit-3@YST`. List the pairs in input order, the way the loops above visit them.
