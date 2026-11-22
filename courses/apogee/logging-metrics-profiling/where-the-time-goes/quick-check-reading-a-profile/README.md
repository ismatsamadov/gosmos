---
slug: quick-check-reading-a-profile
title: Quick Check: Reading a Profile
type: quiz
position: 1
difficulty: 8
topics: profiling
created: 2026-11-10T15:36:39Z
modified: 2026-11-11T17:33:38Z
---
Tuesday afternoon. Nkem ran `go tool pprof -top` on a profile of Tuesday's plan, made on the planning machine with your `profileTo`, and copies the top six lines onto the whiteboard, counting everything `elevationDeg` calls, its sines and cosines included, as its own.

`flat` is time in a function's own code; `cum` adds all it called:

```text
  flat flat%    cum   cum%
188.0s 81.0% 188.0s  81.0% elevationDeg
 12.4s  5.3%  12.4s   5.3% scanAllBookings
 10.2s  4.4% 199.6s  86.0% scanEverySecond
  8.9s  3.8%  19.9s   8.6% sortAfterEach
  1.4s  0.6% 189.4s  81.6% elevation.func1
  0.1s  0.0% 232.0s 100.0% planDay
```

## Your task

Answer Nkem's three questions: whose own code takes the most time, what `planDay`'s line says, and where the first gain is.
