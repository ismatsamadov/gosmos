---
slug: regions-for-each-step
title: Regions for Each Step
type: code_tests
position: 2
difficulty: 9
topics: tracing
created: 2026-11-10T15:36:37Z
modified: 2026-11-10T15:36:37Z
---
Thursday 15 February. Since Christmas, Tomás had blamed a hung radio: "The radio hung, so the worker never came back." For the review, a trace must show where a failed pass spends its time: waiting for a worker, or flying.

A region marks one step of a task, in one goroutine. `trace.WithRegion(ctx, name, f)` runs `f` as a region; `trace.StartRegion(ctx, name)` returns a region you end yourself with `End`, in the same goroutine, for a step that isn't one function:

```go
r := trace.StartRegion(ctx, "read settings")
cfg, err := load()
r.End()
trace.WithRegion(ctx, "plan", func() {
	plan(cfg)
})
```

## Your task

In `runPass`, put the wait for a worker, the `select`, in a region named `wait for worker`, ended on both of its ways out, and the call to `fly` in a region named `fly`.

> [!NOTE]
> Regions change nothing a pass returns, so the test also reads `uplink.go` and checks that every `StartRegion` is ended in `runPass`.
