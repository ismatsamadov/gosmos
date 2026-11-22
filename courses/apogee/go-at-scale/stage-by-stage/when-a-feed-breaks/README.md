---
slug: when-a-feed-breaks
title: When a Feed Breaks
type: code_tests
position: 6
difficulty: 8
topics: context, errors
created: 2026-11-10T15:39:51Z
modified: 2026-11-11T12:51:41Z
---
In the drill, Ystmark's link to Port Ellery dropped in the middle of its feed, and the old loop carried on as if Ystmark's satellites had gone quiet.

`context.WithCancelCause` returns a cancel function that takes an error, the cause. The given `readFeed` calls it when its reader fails, with the error wrapped as `feed YST: ...`, and every stage stops. Then `ctx.Err()` only says `context canceled`, while `context.Cause(ctx)` says why.

```go
ctx, cancel := context.WithCancelCause(parent)
defer cancel(nil) // later calls change nothing
cancel(errors.New("wind 74 km/h"))
fmt.Println(context.Cause(ctx)) // wind 74 km/h
```

## Your task

`ingest` cancels the run without a cause, so a broken feed goes unreported. Make its context with `context.WithCancelCause`, hand that cancel to every `readFeed`, and once `newest` returns, return `context.Cause(ctx)` when it isn't `nil`.
