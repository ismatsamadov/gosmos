---
slug: five-stations-one-moment
title: Five Stations, One Moment
type: code
position: 0
difficulty: 6
topics: goroutines
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:26Z
---
Professor Marr, who leads Thistle's student team, wants every one of the satellite's last passes heard. The telemetry service polls stations one after another: by the time it reaches Lantern Atoll, the pass is over.

`go f()` starts `f` in a new goroutine, running alongside the rest of the program, and doesn't wait for it. When `main` returns, the program ends, whatever its goroutines are still doing. A `sync.WaitGroup` does the waiting: `wg.Go(f)` starts `f` in a goroutine and counts it, and `wg.Wait()` blocks until all of them have returned.

```go
var wg sync.WaitGroup
dishes := []string{"La Grande", "Chiquita"}
for _, dish := range dishes {
	wg.Go(func() { calibrate(dish) })
}
wg.Wait()
fmt.Println("both dishes calibrated")
```

## Your task

`main` starts the polls with `go` and prints before any has answered. Start them with `wg.Go` instead, and call `wg.Wait()` before printing. Add `"sync"` to the imports.

> [!NOTE]
> Each goroutine writes only its own `counts[i]`, so they never collide.
