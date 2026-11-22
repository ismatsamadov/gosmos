---
slug: a-task-for-each-pass
title: A Task for Each Pass
type: code_tests
position: 1
difficulty: 9
topics: context, defer, tracing
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:39Z
---
Wednesday 14 February. In `go tool trace`, the bench replay is thousands of goroutine switches, and nothing says which pass any of them worked for. Nkem: "Which of these were Cirrine-1's? Every trace from now on should say."

A task groups a trace's events under a name. `trace.NewTask(ctx, name)` starts one and returns a context that carries it: whatever is recorded with that context, in any goroutine, belongs to the task, until `task.End()`. `trace.Log(ctx, category, message)` records a note, such as whose task it is:

```go
ctx, task := trace.NewTask(ctx, "plan day")
defer task.End()
trace.Log(ctx, "station", "LAN")
planDay(ctx)
```

## Your task

At the top of `runPass`, start a task named `pass` and end it with `defer`. Log the pass ID with category `pass`, and hand the task's context on to the `select` and to `fly`.

> [!NOTE]
> A task changes nothing a pass returns, so the test also reads `uplink.go` with `go/parser` and says what it looked for.
