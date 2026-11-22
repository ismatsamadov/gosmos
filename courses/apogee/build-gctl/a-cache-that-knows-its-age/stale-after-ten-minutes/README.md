---
slug: stale-after-ten-minutes
title: Stale After Ten Minutes
type: code_tests
position: 1
difficulty: 7
topics: functions, time
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T17:14:45Z
---
Ingrid has one condition for the cache: "A schedule changed at 04:00 reaches my prompt by 04:10. Copy?" Today the cache keeps its first load all night.

The cache reads the time through its `now` field instead of calling `time.Now` itself, so a test can give it any time at all and never wait. Comparing times takes `Add` and `Before`:

```go
los := aos.Add(10 * time.Minute)
inView := now.Before(los)
```

## Your task

1. Write a `refresh` method that calls `c.load`, keeps the passes, and records `c.now()` in `c.loadedAt`.
2. Make `passes` call `refresh` when nothing is loaded or the passes are stale: when `c.now()` is no longer before `c.loadedAt` plus `c.ttl`.
3. In `newShell`, give the cache `time.Now` as its clock and a time to live of 10 minutes.
