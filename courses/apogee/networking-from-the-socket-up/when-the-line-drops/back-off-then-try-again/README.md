---
slug: back-off-then-try-again
title: Back Off, Then Try Again
type: code_tests
position: 1
difficulty: 8
topics: networking, time
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Wednesday. Sunniva has done the sums: a dial every second for fourteen minutes is 840 dials, on a line the station needs for everything else. "Wind 30, easing. 840 dials. Too many."

Back off instead: wait longer after each failure in a row, up to a cap, and start from the shortest wait again once a dial works. The given `redialDelay(attempt, r)` gives 1 s after the first failure, doubles for each one after it, stops at 30 s, and multiplies by a jitter, a random factor from 0.5 to 1 drawn from `r`:

```go
r := rand.New(rand.NewPCG(1, 2))
fmt.Println(redialDelay(1, r), redialDelay(2, r),
	redialDelay(3, r))
// 838.227829ms 1.461386217s 3.017094795s
```

## Your task

`redial` waits one second every time. Count failed dials in `l.attempt`, wait `redialDelay(l.attempt, l.Rand)` after each, and set `l.attempt` back to 0 when a dial works.

> [!TIP]
> Count a failure before its wait: the first failed dial is attempt 1.
