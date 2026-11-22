---
slug: one-stations-day-in-order
title: One Station's Day, in Order
type: code_tests
position: 7
difficulty: 9
topics: generics, iterators
created: 2026-11-11T12:51:43Z
modified: 2026-11-11T12:51:43Z
---
Friday 13 April. Cirrine's twelve new satellites leave their rocket at 09:40 on Tuesday, and Tirrin Rock's first day with them must be planned from eighty streams, without one big sort.

The pieces fit together. `Passes` predicts one satellite's passes a day at a time, `Filter` keeps the ones that climb high enough, and `MergeBy` merges the streams in `byAOS` order. Building the streams predicts nothing: each is pulled only as far as its head, once a loop asks. In `types.go`, a `Network`'s `Satellites` and `Elements` list the eighty in one order:

```go
for i, sat := range n.Satellites {
	e := n.Elements[i] // the same satellite
	fmt.Println(sat.Name, e.Rev0)
}
```

## Your task

`stationPasses` yields nothing. For each satellite `s` can serve (`canServe`), filter its `Passes` over `w` to those whose max elevation is at least `minElevDeg`, then return all of them merged by `MergeBy` and `byAOS`.

> [!NOTE]
> The tests hand Tuesday's 148 candidates to `assign`, the planning rule in `plan.go`, which gives 120 of them one of Tirrin Rock's two dishes: 120 contacts.
