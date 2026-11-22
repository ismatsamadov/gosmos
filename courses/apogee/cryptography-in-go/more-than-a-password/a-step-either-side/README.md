---
slug: a-step-either-side
title: A Step Either Side
type: code_tests
position: 4
difficulty: 9
topics: crypto, loops, time
created: 2026-11-13T06:18:24Z
modified: 2026-11-13T06:18:24Z
---
Tuesday 31 October. Ingrid tries the test console and is refused twice with the code her phone shows. The phone's clock runs 20 seconds behind the console's, so for 20 seconds of every 30 it still shows the code of the step before.

Two clocks never agree exactly, and reading and typing a code takes seconds. So the console also accepts the codes of the step before and the step after, and no more: each step it accepts is one more code a guess could match. Steps are plain counters, and `hotp` makes any step's code:

```go
s := timeStep(time.Unix(89, 0))
fmt.Println(s-1, s, s+1) // 1 2 3
```

## Your task

Make `Check` accept `code` when it is the code of `user`'s secret for the step at `c.now()`, the step before or the step after, still compared in constant time.

> [!NOTE]
> The console keeps each secret itself, never a digest: it must make the same codes as the phone, so it guards every secret like a key.
