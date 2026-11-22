---
slug: pass-day
title: Pass Day
type: code_tests
position: 8
difficulty: 9
topics: context, errors, goroutines, mutexes
created: 2026-11-01T07:36:47Z
modified: 2026-11-06T16:51:27Z
---
Friday, your last on-call shift: a full day of passes, and a storm building at Ystmark. "Wind 64, rising. Expect a stow," Sunniva writes.

Concurrent code is easiest to trust when goroutines share little. Let each goroutine do its slow work alone, and hold the lock only to add its result. Without the lock, two stations that finish together could update the total at the same moment and lose a count, and no test run is sure to catch it.

```go
for _, code := range codes {
	wg.Go(func() {
		n := countFrames(code) // slow, no lock
		mu.Lock()
		total += n
		mu.Unlock()
	})
}
```

## Your task

1. In `runDay`, run each station's passes in its own goroutine, add each station's report to the total under a mutex, and wait for all of them.
2. In `runPass`, give `run` a context that ends at the pass's LOS.
3. In `runStation`, once a pass fails with `ErrStowed`, count the station's remaining passes as skipped and stop.

> [!TIP]
> `run` wraps `ErrStowed` with details, so look for it with `errors.Is`.
