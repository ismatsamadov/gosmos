---
slug: quick-check-pipelines
title: Quick Check: Pipelines
type: quiz
position: 7
difficulty: 7
topics: channels, context
created: 2026-11-10T15:39:51Z
modified: 2026-11-11T12:51:41Z
---
Thursday 8 March, 10:20: Thistle-2, Kessick Polytechnic's new student satellite, leaves its rocket, and its first beacon is due over Port Ellery at 17:08. Before your pipeline carries it to the Board, Nkem has three questions from the drill.

The first starts from a send like this one:

```go
select {
case out <- f:
case <-ctx.Done():
	return
}
```

## Your task

Answer Nkem's questions: what happens to a stage whose send has no `ctx.Done()` case once the Board stops reading, what `ctx.Err()` and `context.Cause(ctx)` say after `readFeed` cancels the run, and when a fan-in may close its output.
