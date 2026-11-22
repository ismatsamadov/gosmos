---
slug: one-page-per-problem
title: One Page per Problem
type: code_tests
position: 4
difficulty: 8
topics: methods, metrics, time
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
The old pager sent a page at every LOS while a problem lasted, and so far the new one is worse: a FIRING every minute. Tomás: "Eleven times. Like a buoy bell in a swell, and the same problem every time."

A problem is one rule at one station, from the FIRING that starts it to the RESOLVED that ends it. Send FIRING once; while the rule stays firing, send it again only when `Repeat` has passed since the last one, and never when `Repeat` is 0. Only a Critical rule sends a RESOLVED: a Warning that clears has woken nobody.

```go
repeat := 24 * time.Hour
since := 23*time.Hour + 59*time.Minute
fmt.Println(since >= repeat) // false: wait
```

## Your task

`whileFiring` sends FIRING at every sample. Send nothing while the rule stays firing, unless `r.Repeat` isn't 0 and at least `r.Repeat` has passed since `st.sent`: then set `st.sent` to `s.At` and send FIRING. Send RESOLVED only when `r.Severity` is `Critical`.

> [!NOTE]
> `beforeFiring` sets `st.sent` when a rule fires, so the first reminder counts from the first FIRING.
