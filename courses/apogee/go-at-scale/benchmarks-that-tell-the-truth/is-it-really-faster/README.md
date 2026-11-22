---
slug: is-it-really-faster
title: Is It Really Faster?
type: code_tests
position: 7
difficulty: 9
topics: performance, testing
created: 2026-11-11T22:45:16Z
modified: 2026-11-12T08:43:36Z
---
Friday 18 May. Tomás's forty comes off the glass today, and Nkem wants the truth in its place: the old planner and passplan 2 on the same network and windows, six hours, half a day and a day.

`b.Run(name, f)` runs `f` as a sub-benchmark, with its own loop and its own line, named after its parents, as in `BenchmarkWindows/old/hours=6`. Sub-benchmarks nest, so one run can line up every version at every size:

```go
for _, size := range []int{10, 100, 1000} {
	name := fmt.Sprintf("lines=%d", size)
	b.Run(name, func(b *testing.B) {
		for b.Loop() {
			parseAll(lines[:size])
		}
	})
}
```

## Your task

`benchWindows` measures `planNew` on a day. Give it a sub-benchmark `old` for `planOld` and one `new` for `planNew`, each with sub-benchmarks `hours=6`, `hours=12` and `hours=24` that plan that many hours from midnight on 18 May, with the same network and settings, and report their allocations.

> [!NOTE]
> In a terminal, `go test -bench=Windows -count=10` lines them all up. A sub-benchmark doesn't inherit `b.ReportAllocs` from its parent: call it where the loop is.
