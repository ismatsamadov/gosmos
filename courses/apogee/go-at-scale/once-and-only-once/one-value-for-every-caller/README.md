---
slug: one-value-for-every-caller
title: One Value for Every Caller
type: code_tests
position: 1
difficulty: 8
topics: functions, goroutines
created: 2026-11-11T12:51:42Z
modified: 2026-11-11T12:51:42Z
---
Monday afternoon, Nkem reads your `Tables` and draws a circle around its three fields: "A `Once` and two fields for one table. Go has a function for that."

`sync.OnceValues(f)` returns a function that calls `f` the first time it is called, from whichever goroutine, and returns `f`'s two results on every call after that without calling `f` again. Callers that arrive while `f` runs wait for it, as with `once.Do`. It keeps an error as firmly as a value. For a function with one result, there is `sync.OnceValue`.

```go
gust := sync.OnceValues(func() (int, error) {
	return readMast("CWR")
})
kmh, err := gust() // reads the mast
kmh, err = gust()  // the same two, no read
```

## Your task

`networkOnce` returns `load` itself, so every call loads again. Return `sync.OnceValues(load)` instead.

> [!WARNING]
> If the file server is down at the first call, every later call gets that error too. A load that should try again needs a mutex and a check of its own, not `OnceValues`.
