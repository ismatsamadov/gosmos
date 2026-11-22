---
slug: first-failure-cancels
title: First Failure Cancels the Rest
type: code_tests
position: 7
difficulty: 8
topics: context, errors, goroutines
created: 2026-11-01T07:36:47Z
modified: 2026-11-06T16:51:27Z
---
Before Furrow-A's evening passes, every station that will track it runs its pre-pass check at once. If Ystmark's fails, the plan is off and the other checks are wasted.

`context.WithCancelCause(ctx)` returns a child context and a cancel function that takes an error, the cause. The first call cancels the context and records its cause; later calls change nothing. `context.Cause(ctx)` gives the cause back, or `nil` while the context is still running.

```go
parent := context.Background()
ctx, cancel := context.WithCancelCause(parent)
cancel(errors.New("storm warning: stow dishes"))
cancel(errors.New("too late, ignored"))
fmt.Println(ctx.Err()) // context canceled
cause := context.Cause(ctx)
fmt.Println(cause) // storm warning: stow dishes
```

## Your task

`runAll` ignores what the checks return. Give them a child context from `context.WithCancelCause`, cancel it with the error of any check that fails, and once all have returned, return the cause.
