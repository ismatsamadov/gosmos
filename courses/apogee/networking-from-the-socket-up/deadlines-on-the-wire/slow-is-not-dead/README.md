---
slug: slow-is-not-dead
title: Slow Is Not Dead
type: code_tests
position: 3
difficulty: 8
topics: interfaces, networking, time
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Wednesday, the fiber still down: dishd gives each frame 5 seconds, and over the backhaul each 27-byte STATUS from Dish 1 takes ten, so every one is cut.

A slow line keeps bringing bytes; a dead one brings none. An idle deadline tells them apart: it moves forward before every `Read`, so each `Read` must bring a byte within the timeout. A reader wrapped around the connection can act before every read `c.br` makes:

```go
type logReader struct{ r io.Reader }

func (l logReader) Read(p []byte) (int, error) {
	log.Println("read") // before every read
	return l.r.Read(p)
}
```

## Your task

1. In `idleReader`'s `Read`, set the read deadline of `r.c.conn` to now plus `r.c.Timeout`, then read.
2. Take the deadline out of `recv`: `newClient` now makes `c.br` read through `idleReader`.

> [!TIP]
> A `bufio.Reader` calls `Read` only when its buffer runs dry: a frame arriving whole costs one deadline, one that trickles in gets one per piece.
