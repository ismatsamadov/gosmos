---
slug: deliver-with-retries
title: Deliver With Retries
type: code_tests
position: 6
difficulty: 9
topics: context, errors, http, loops
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Thursday: Cirrine's receiver is in maintenance again and answers `503` twice before it is back. Its event has to arrive on the third try without anyone at Apogee lifting a finger.

So `Deliver` loops. Each attempt is signed with its own time from `s.Now()`. It ends on a 2xx; it ends with `ErrUndelivered` on a failure `retryable` won't retry, or after the 8th attempt; otherwise it waits and tries again. The wait goes through `s.Wait`, which sleeps in production but only takes notes in tests, so eight attempts take no time:

```go
s.Wait = func(ctx context.Context,
	d time.Duration) error {
	waits = append(waits, d)
	return ctx.Err()
}
```

## Your task

`Deliver` tries once. Make it try again while `retryable(statusOf(err), err)` allows, up to 8 attempts in all, after `s.Wait(ctx, s.schedule(d))`. When `ctx` ends, return `ctx.Err()` and keep the count in `d.Attempts`: a service that stops hasn't failed.
