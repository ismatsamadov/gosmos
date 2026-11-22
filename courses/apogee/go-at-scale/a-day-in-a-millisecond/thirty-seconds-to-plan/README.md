---
slug: thirty-seconds-to-plan
title: Thirty Seconds to Plan
type: code_tests
position: 4
difficulty: 8
topics: context, errors, time
created: 2026-11-11T16:51:57Z
modified: 2026-11-12T08:43:36Z
---
Thursday 26 April. A planning run that hangs holds up every replan behind it, so Nkem wants each run over within thirty seconds, with an error that says why.

`context.WithTimeout(ctx, d)` returns a child context that ends after `d`, or sooner if its parent ends, and a `cancel` function to call once the work is done, which stops the child's timer. When time runs out, the child's `Err()` is `context.DeadlineExceeded`, and `%w` keeps that error findable by `errors.Is` inside a clearer one:

```go
ctx, cancel := context.WithTimeout(ctx, wait)
defer cancel()
err := upload(ctx, "OST") // stops after wait
if errors.Is(err, context.DeadlineExceeded) {
	err = fmt.Errorf("upload: too slow: %w", err)
}
```

## Your task

`planWithin` calls `plan(ctx)` as it is. Give the run a context that ends after `limit`, cancel it once the run returns, and when the run's error is a deadline, wrap it with `%w` after this message, which names `limit`:

```text
plan: took longer than 30s
```

> [!NOTE]
> A run that never looks at its context runs on: the timeout only closes `ctx.Done()`.
