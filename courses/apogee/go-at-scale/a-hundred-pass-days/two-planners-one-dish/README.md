---
slug: two-planners-one-dish
title: Two Planners, One Dish
type: code_tests
position: 3
difficulty: 9
topics: maps, slices, time
created: 2026-11-12T03:18:38Z
modified: 2026-11-12T10:15:10Z
---
Monday 25 June, the rehearsal. passplan 2 planned the day from scratch at midnight, and six of its first contacts took a dish the old planner was still using.

Nkem's question again: what happens when two planners touch the same pass? At the cut-over, the instant the new planner takes over, the old one keeps flying its contacts in progress. A dish is free after its last contact's LOS and the turnaround, two minutes to turn and settle: the hand-over gives `Plan` those times as `free`. A slice in a map can grow to fit an index:

```go
seats := map[string][]bool{}
for _, b := range bookings {
	s := seats[b.Car]
	for len(s) < b.Seat {
		s = append(s, false)
	}
	s[b.Seat-1], seats[b.Car] = true, s
}
```

## Your task

`handover` returns nil. For each contact of `cur` that rose before `at`, keep the latest LOS plus `turnaround` of its station's dish `d`, at index `d-1`. Contacts come in AOS order: stop at the first that rises at `at` or later.

> [!NOTE]
> At 00:00 on 28 June, fourteen dishes are still busy: twelve with contacts in progress, two still settling.
