---
slug: planned-in-seconds
title: Planned in Seconds
type: code_tests
position: 7
difficulty: 9
topics: maps, profiling, sorting
created: 2026-11-10T15:36:39Z
modified: 2026-11-11T17:33:38Z
---
Friday 19 January. At 14:00 the new passplan goes out, and Nkem has asked for one more thing in that deploy: every program serves its profiles on its debug port, the private address where the night shift already reads its numbers, for the next slow thing.

The three fixes only add up when `planDay` uses them together: coarse steps for each satellite and station, then the bookings and the order for the whole day at once, so the index is built and the day sorted one time, not once per station:

```text
scanEverySecond  ->  findPasses, per station
scanAllBookings  ->  attach, once a day
sortAfterEach    ->  buildDay, once a day
```

## Your task

1. Make `planDay` find each satellite's passes over each station with `findPasses`, and gather them all in one slice.
2. Then call `attach` once on them, and `buildDay` once on what it returns.

> [!TIP]
> `go tool pprof -top cpu-after.prof` shows where the time goes now: the test profiles your planner into that file with your `profileTo`.
