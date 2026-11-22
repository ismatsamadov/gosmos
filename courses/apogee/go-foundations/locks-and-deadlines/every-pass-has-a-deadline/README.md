---
slug: every-pass-has-a-deadline
title: Every Pass Has a Deadline
type: code_tests
position: 4
difficulty: 7
topics: context, time
created: 2026-11-01T07:36:47Z
modified: 2026-11-06T16:51:27Z
---
A pass ends at LOS whether the code is finished or not. Ingrid wants each pass's context to end there on its own, with nobody pressing cancel.

`context.WithDeadline(parent, t)` returns a context that ends at time `t`, or earlier if `parent` ends first. `Deadline()` reports `t` (or the parent's deadline, if that is sooner), and when the deadline ends the context, `Err()` returns `context.DeadlineExceeded`. Call the cancel function it returns even though the deadline ends the context on its own: `defer cancel()` releases it as soon as the work is done.

```go
stowBy := time.Date(2027, 6, 29,
	2, 0, 0, 0, time.UTC)
ctx, cancel := context.WithDeadline(
	context.Background(), stowBy)
defer cancel()
// gives up at 02:00 UTC
err := stowDishes(ctx, "YST")
```

## Your task

`passContext` returns a context that never ends on its own. Make it end at `los`, keeping `parent` as its parent.
