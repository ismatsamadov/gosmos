---
slug: storm-at-ystmark
title: Storm at Ystmark
type: code
position: 9
difficulty: 6
topics: loops, methods, pointers, slices
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:26Z
---
The storm reaches Ystmark: 64 km/h at 02:00, then 75, 58 and 81 by 05:00. Both dishes should stow from 72 km/h and track again below it, but your log says they never moved.

In `for _, d := range dishes`, `d` is a copy of each element, so a pointer method called on `d` changes the copy. Index the slice to reach the element itself:

```go
for _, st := range stations {
	st.Stow("storm") // stows a copy
}
for i := range stations {
	stations[i].Stow("storm") // stows the element
}
```

Reading through a copy is fine, as `state` does; changing one isn't.

## Your task

Fix the loop in `react` so that it stows and unstows the dishes in the slice. The program should print:

```text
02:00 wind 64 km/h: tracking
03:00 wind 75 km/h: stowed
04:00 wind 58 km/h: tracking
05:00 wind 81 km/h: stowed
dish 1 stowed 2 times
dish 2 stowed 2 times
```
