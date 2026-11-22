---
slug: the-shadow-run
title: The Shadow Run
type: code_tests
position: 8
difficulty: 9
topics: iterators, slices, time
created: 2026-11-11T22:52:57Z
modified: 2026-11-11T22:52:57Z
---
Thursday 5 April. The shadow run compares a whole day for the first time, and at Tirrin Rock the planners disagree: its host parked both dishes from 07:30 to 08:24 to clear their ice, and only the old planner was told.

Both schedules come in `byAOS` order, which leaves contacts that rise in the same second in any order. At 05:00:11 Brine-3 rose over Salar Quieto and Furrow-A over Varsk Island, and the two planners list them differently. So walk the schedules a second at a time: gather what each has in that second, then match those. A pull loop can take a run of values:

```go
for ; ok && g.Kmh >= 72; g, ok = next() {
	gale = append(gale, g)
}
```

## Your task

`Compare` yields nothing. Pull from both schedules, with both stopped on every way out, and until both have ended:

1. take the earliest AOS either has left;
2. collect each schedule's contacts with that AOS;
3. yield what `changes` (`match.go`) finds among them, returning as soon as `yield` returns false.
