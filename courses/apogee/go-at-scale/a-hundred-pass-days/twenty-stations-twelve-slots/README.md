---
slug: twenty-stations-twelve-slots
title: Twenty Stations, Twelve Slots
type: code_tests
position: 1
difficulty: 9
topics: context, goroutines, iterators
created: 2026-11-12T03:18:38Z
modified: 2026-11-12T03:18:38Z
---
Tuesday 12 June. Planned one after another, the twenty stations keep one of passplan 2's twelve cores busy, and Nkem wants all twelve at work.

No two stations share a dish, so stations can be planned at once, in any order, and the plan is the same. Run each station as a job of chapter 2's `Group`; a job that writes only its own index of a slice needs no lock:

```go
g, _ := newGroup(ctx, 2)
sizes := make([]int, len(files))
for i, f := range files {
	g.Go(func(ctx context.Context) (err error) {
		sizes[i], err = size(ctx, f)
		return err
	})
}
```

Each station's list is in AOS order already, so chapter 6's `MergeBy` can join the twenty without a sort.

## Your task

`Plan` plans the stations one after another, then sorts every contact. Plan each station as a job of a group of `cfg.Slots`, return the error `Wait` returns, and merge the lists with `MergeBy` and `byContact`.

> [!NOTE]
> The tests hold twelve stations inside the pass model: no thirteenth may start, and after a cancel no other station may.
