---
slug: read-over-a-shoulder
title: Read Over a Shoulder
type: code_tests
position: 5
difficulty: 9
topics: crypto, maps, security
created: 2026-11-13T06:18:24Z
modified: 2026-11-13T06:18:24Z
---
Wednesday 1 November. Mirela watches Nkem type a code into the test console, types the same six digits on her own laptop 40 seconds later, and the code check accepts them again. "Anything I can read over a shoulder must work once."

RFC 6238 says so too: once a code is accepted, it never counts again, even inside its window. `Codes` now keeps `last`, the step of the last code accepted for each user, and a code counts only when its step is above that one, the way the replay guard treats sequence numbers. A map gives 0 for a user it doesn't hold yet:

```go
last := map[string]uint64{}
fmt.Println(last["nkem"]) // 0
last["nkem"] = 61889550
```

## Your task

Make `Check` skip every step that isn't above `c.last[user]`, and record in `c.last` the step of the code it accepts. A refused code changes nothing.
