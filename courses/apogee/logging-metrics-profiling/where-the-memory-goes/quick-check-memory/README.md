---
slug: quick-check-memory
title: Quick Check: Memory
type: quiz
position: 3
difficulty: 9
topics: goroutines, profiling
created: 2026-11-10T15:36:39Z
modified: 2026-11-11T17:33:38Z
---
Wednesday 7 February, 09:00. At 02:00 you saved Lantern Atoll's profiles from uplinkd's debug port, and Nkem reads them with you, goroutines first.

Line 13 of `uplink.go` is `u.workers <- struct{}{}`, and two lines below it, line 15 is `done <- result{rep, err}`; the other 14 goroutines are uplinkd's own. `heap.prof` holds 29 MiB in use, 20 MiB of it five buffers made in `runPass`.

```text
goroutine profile: total 19
4 @ 0x46f241 0x40a2b4 0x6c41d6
#	main.(*Uplink).runPass.func1	uplink.go:15
1 @ 0x46f241 0x40a2b4 0x6c4107
#	main.(*Uplink).runPass.func1	uplink.go:13
```

## Your task

Answer Nkem's three questions: why `heapProfile` collects garbage before it writes, what the four goroutines on line 15 are doing, and why each keeps 4 MiB in use.
