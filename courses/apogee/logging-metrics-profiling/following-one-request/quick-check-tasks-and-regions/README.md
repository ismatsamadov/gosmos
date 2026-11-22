---
slug: quick-check-tasks-and-regions
title: Quick Check: Tasks and Regions
type: quiz
position: 3
difficulty: 9
topics: goroutines, tracing
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:39Z
---
Thursday afternoon. Over the bench trace, Nkem sketches a pass whose flight runs on a goroutine of its own, as the old uplinkd's did, and asks you three questions about it.

On the whiteboard:

```go
ctx, task := trace.NewTask(ctx, "pass")
defer task.End()
go func() {
	r := trace.StartRegion(ctx, "fly")
	defer r.End()
	fly(ctx)
}()
```

## Your task

Answer Nkem's three questions: which goroutines the task covers, where the region must end, and what the task and the region cost when no trace is running.
