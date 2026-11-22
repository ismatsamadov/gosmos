---
slug: hang-up-when-done
title: Hang Up When Done
type: code_tests
position: 1
difficulty: 6
topics: defer, errors, networking
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T07:44:27Z
---
Monday evening, after an afternoon of checks, the simulator's log shows forty connections that dishd, the dish service you are writing, opened and never closed. A controller serves one client at a time, so it threw each one out as the next arrived.

A connection stays open at both ends until one of them closes it. `defer conn.Close()` closes it when the function returns, whichever `return` that is. Put it right after the error check:

```go
f, err := os.Open("manual.txt")
if err != nil {
	return err
}
defer f.Close() // runs when the function returns
```

## Your task

`bannerOf` dials a controller and returns the line every controller sends first, its banner, read by the given `readBannerLine`. It leaves the connection open: close it with a `defer` in the right place.

> [!WARNING]
> A `defer conn.Close()` above the error check panics when nobody is listening: a failed dial leaves `conn` nil.
