---
slug: at-this-rate
title: At This Rate
type: code_tests
position: 5
difficulty: 8
topics: metrics, time
created: 2026-11-10T15:36:37Z
modified: 2026-11-10T15:36:37Z
---
Wednesday 10 January, from Yusra: "TL;DR: at this rate, when does each station run out?" Customers want a day they can plan for, not a share.

At the last week's rate a station loses `lostLast7d` passes every 7 days, so what remains of its budget lasts remaining × 7 days ÷ `lostLast7d`. A `time.Duration` counts nanoseconds in an `int64`, so multiply before you divide, or the answer rounds to whole weeks:

```go
week := 7 * 24 * time.Hour
fmt.Println(3*week/2, 3/2*week)
// 252h0m0s 168h0m0s
```

## Your task

Make `exhaustIn` return how long `remaining` passes of budget last at `lostLast7d` lost a week, and true. Nothing left, or a budget overspent, gives 0 and true, whatever the week lost; nothing lost this week gives 0 and false, since there is no rate to go by.

> [!NOTE]
> A station already past its budget, as Lantern Atoll is, has run out: its answer is 0, now.
