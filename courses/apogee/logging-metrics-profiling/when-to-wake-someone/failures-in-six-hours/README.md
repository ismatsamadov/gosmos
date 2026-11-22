---
slug: failures-in-six-hours
title: Failures in Six Hours
type: code_tests
position: 0
difficulty: 8
topics: loops, metrics, time
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
Ingrid is rewriting the pager's rules, and failed passes come first: "One problem, one page. And every rule starts from one number."

For failures that number is a count over a window of time: a station's failed passes whose LOS, the moment the satellite set, falls in the last six hours. The window holds the times after `now` minus six hours and up to `now` itself, so at 01:41:28 it no longer holds the night's first failure, at 19:41:28:

```go
now := time.Date(2028, 12, 28, 1, 41, 28, 0,
	time.UTC)
from := now.Add(-6 * time.Hour)
fmt.Println(from.Format("Jan 2 15:04:05"))
// Dec 27 19:41:28
```

## Your task

`failedWithin` returns 0. Count the results for `station` whose `Outcome` is `failed` and whose `LOS` is after `now` minus `d` and not after `now`.

> [!NOTE]
> A skipped pass, one the station never tried to fly, isn't a failure: only passes that went wrong count.
