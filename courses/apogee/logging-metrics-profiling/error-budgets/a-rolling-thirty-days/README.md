---
slug: a-rolling-thirty-days
title: A Rolling Thirty Days
type: code_tests
position: 3
difficulty: 8
topics: metrics, time
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:38Z
---
Tuesday 9 January, the 08:00 meeting. Ingrid: "Customers don't wait for a quarter to end. From tomorrow, each morning starts with every station's last 30 days."

A rolling window moves with the clock: each morning it covers the 30 days before now. Written as the half-open range (now − 30 days, now], it counts a pass whose LOS, the moment it set, is exactly now, and leaves out one that set exactly 30 days before, so two windows that meet never count a pass twice. `now` comes in as an argument, as in the pager's rules, so a test picks the morning:

```go
since := now.Add(-6 * time.Hour)
in := func(los time.Time) bool {
	return los.After(since) && !los.After(now)
}
fmt.Println(in(since), in(now))
// false true
```

## Your task

Make `window` count only `station`'s results whose LOS is in (now − 30 days, now]: `booked` counts each of them, and `lost` those that failed or were skipped.

> [!NOTE]
> A skipped pass is lost too: a dish parked out of the wind, stowed, still leaves a customer without the pass.
