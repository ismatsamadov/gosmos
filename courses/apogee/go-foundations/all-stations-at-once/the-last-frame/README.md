---
slug: the-last-frame
title: The Last Frame
type: code
position: 3
difficulty: 6
topics: channels
created: 2026-11-01T07:36:47Z
modified: 2026-11-01T07:36:47Z
---
Ystmark's archiver can't tell when Pipit-1's pass is over, so it waits forever for a frame that never comes.

A sender with nothing more to send calls `close(ch)`. A `for v := range ch` loop receives until the channel is closed and empty, then ends. After that, every receive returns the zero value at once, and `v, ok := <-ch` sets `ok` to `false`.

```go
pongs := make(chan string, 2)
pongs <- "PEL"
pongs <- "SQT"
close(pongs)
for code := range pongs {
	fmt.Println("pong from", code)
}
```

## Your task

Run the program: it archives four frames, then dies with `all goroutines are asleep - deadlock!`. In `poller`, close the channel after the last send.

> [!WARNING]
> Only the sender closes, and only once: closing a closed channel panics, and so does sending on one.
