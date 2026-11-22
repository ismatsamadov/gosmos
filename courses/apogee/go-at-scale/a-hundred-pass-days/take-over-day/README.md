---
slug: take-over-day
title: Take-Over Day
type: code_tests
position: 8
difficulty: 10
topics: context, goroutines, select, time
created: 2026-11-12T03:18:38Z
modified: 2026-11-12T10:15:10Z
---
Thursday 28 June 2029, 00:00:00 UTC, the cut-over. Twenty-four new Furrows were heard for the first time last night, eighty satellites fly over twenty stations, and Cape Wren's wind is due at 09:47.

Run 1 is a replan too, of the old planner's schedule from the cut-over: `replanFrom` keeps its contacts, which the old planner flies, and plans around the dishes they hold. Every later replan then knows those dishes; without them, a replan at 00:10 would give Ilvane Spit's second dish to a Cirrine satellite three seconds before it settles. The station runner flies the newest schedule, waking for the next contact's AOS or for news:

```go
select {
case <-time.After(time.Until(next)):
	ring()
case <-moved: // the timetable changed
case <-ctx.Done():
	return ctx.Err()
}
```

## Your task

1. Give a copy of `old` passplan 2's window, from `at` for `cfg.Days` days, and replan it from `at`; publish that as version 1.
2. Run a `Replanner` on `events` in a goroutine; wait for it before you return.
3. Until `ctx` ends, fly each contact at its AOS, the first after the last one flown, and read the schedule again when `r.news` has news.
