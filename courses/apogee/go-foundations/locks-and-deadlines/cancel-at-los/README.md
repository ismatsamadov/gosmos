---
slug: cancel-at-los
title: Cancel at LOS
type: code_tests
position: 3
difficulty: 7
topics: context, select
created: 2026-11-01T07:36:47Z
modified: 2026-11-01T07:36:47Z
---
Last night Furrow-B's downlink kept writing after LOS and filled the archive with noise. Its receiver never looked at the context the pass gave it.

A `context.Context` tells a function when to stop. `ctx.Done()` is a channel that closes when the context is canceled, and `ctx.Err()` then says why, such as `context.Canceled`. Put `<-ctx.Done()` in a `select`, next to whatever the function waits for.

```go
select {
case <-ctx.Done():
	return ctx.Err()
case f := <-frames:
	archive(f)
}
```

## Your task

Make `receive` wait with a `select`: when `ctx` is done, return the count so far and `ctx.Err()`; when `chunks` closes, return the count and `nil`.
