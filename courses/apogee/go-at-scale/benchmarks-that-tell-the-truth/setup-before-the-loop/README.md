---
slug: setup-before-the-loop
title: Setup Before the Loop
type: code_tests
position: 1
difficulty: 8
topics: performance, testing
created: 2026-11-11T22:45:16Z
modified: 2026-11-12T08:43:36Z
---
Tuesday 8 May. Tomás's forty came from two benchmarks: the old planner's predicted Cape Wren's day inside its loop and passplan 2's before it, so only one of them timed the prediction.

Everything inside `for b.Loop()` is measured on every turn, in time and in allocations. What comes before the loop runs once and counts for nothing: `b.Loop` starts the clock and the allocation count when the loop starts. So here `AllocsPerOp` is `parseAll`'s alone:

```go
r := testing.Benchmark(func(b *testing.B) {
	lines := drillLines(2000) // once, not counted
	for b.Loop() {
		parseAll(lines)
	}
})
fmt.Println(r.AllocsPerOp())
```

## Your task

`benchPlanDay` now calls `stationDay` inside the loop, as the old planner's benchmark did, so every turn predicts the day again. Move the prediction before `for b.Loop()`, so that only `planDay` is measured.

> [!NOTE]
> The tests count `planDay`'s own allocations with `testing.AllocsPerRun` and allow 2 more per turn; predicting Cape Wren's day makes 2,591.
