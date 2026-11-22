---
slug: many-feeds-one-stream
title: Many Feeds, One Stream
type: code_tests
position: 4
difficulty: 7
topics: channels, goroutines
created: 2026-11-01T07:36:47Z
modified: 2026-11-01T07:36:47Z
---
Four stations hear Thistle at once now, and the archive wants one stream of beacons, not four.

Merging channels, called fan-in, takes one goroutine per input, each sending what it receives into a single output. The output may close only after every input is done, so one more goroutine waits for the others, then closes it. A parameter of type `<-chan int` can only be received from.

```go
var wg sync.WaitGroup
for _, code := range []string{"PEL", "LAN"} {
	wg.Go(func() { replies <- ping(code) })
}
go func() {
	wg.Wait()
	close(replies)
}()
```

## Your task

Write `merge`: give each feed a goroutine that sends everything it receives into `out`. Close `out` once all of them have finished, but return it right away: the caller receives from it.
