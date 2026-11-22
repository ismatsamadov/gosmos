---
slug: wait-longer-each-time
title: Wait Longer Each Time
type: code_tests
position: 4
difficulty: 8
topics: loops, time
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Wednesday: Cirrine's receiver is down for maintenance. A sender retrying every second would knock it over again the moment it came back, with every waiting event at once.

So each wait doubles: 2 seconds after the first failure, 4 after the second, up to 5 minutes. Each is also multiplied by a random factor of at least 0.5 and below 1, called jitter, so senders that failed together don't retry together. A source with a fixed seed draws the same numbers every run:

```go
a := rand.New(rand.NewPCG(1, 2))
b := rand.New(rand.NewPCG(1, 2))
fmt.Println(a.Float64() == b.Float64()) // true
```

## Your task

`backoff` returns 0. Return the wait after failed attempt number `attempt`: 2 seconds doubled once for each earlier failure, at most 5 minutes, times `0.5 + r.Float64()/2`, drawn once.

```text
after attempt 1    1s to 2s
after attempt 2    2s to 4s
after attempt 9    2m30s to 5m0s
```

> [!WARNING]
> `2 * time.Second << 69` overflows to 0: stop doubling at the cap instead.
