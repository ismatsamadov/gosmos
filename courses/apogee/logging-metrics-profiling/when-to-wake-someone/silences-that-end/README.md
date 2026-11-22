---
slug: silences-that-end
title: Silences That End
type: code_tests
position: 6
difficulty: 8
topics: loops, metrics, time
created: 2026-11-10T15:36:39Z
modified: 2026-11-10T15:36:39Z
---
At 07:28 on the 28th Tomás muted Lantern Atoll "until morning", and nothing unmuted it until Ingrid noticed on Friday at 09:00. Her new rule: every silence ends.

A silence holds back one station's alerts from its start, included, to its end, excluded. The given pager now asks `silenced` before it sends an alert, and routes the rest with your `route`. While a silence holds, the rules still follow their values: a problem still firing when it ends is sent at the first sample at or after its end, and one that cleared by then sends nothing. The next silence is planned for Lantern Atoll's generator:

```go
d := time.Date(2029, 1, 11, 0, 0, 0, 0, time.UTC)
gen := Silence{Station: "LAN",
	From:  d.Add(9 * time.Hour),  // included
	Until: d.Add(11 * time.Hour)} // excluded
```

## Your task

`silenced` returns false. Return true when a silence in `s` is for `station` and covers `at`: `at` is not before its `From`, and is before its `Until`.

> [!NOTE]
> `ruleState`, in `state.go`, has a new field, `held`, for a FIRING the silence held back, so the problem pages once the silence is over.
