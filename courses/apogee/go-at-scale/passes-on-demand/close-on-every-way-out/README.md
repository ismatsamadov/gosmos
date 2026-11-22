---
slug: close-on-every-way-out
title: Close on Every Way Out
type: code_tests
position: 6
difficulty: 8
topics: defer, files, iterators
created: 2026-11-11T21:35:30Z
modified: 2026-11-11T21:35:30Z
---
Thursday 29 March, after lunch. A morning of Hollis's questions left the archive server with 412 open files: each loop that found its answer and broke early had left one open.

An iterator's function returns on every way out of the loop over it: after the last value, after an error, and when the loop stops early and `yield` returns false. A `defer` inside it runs then, every time, so it's the place to release what the iterator opened:

```go
return func(yield func(Frame) bool) {
	tune(dish, "S")
	defer park(dish) // the end, an error, a break
	for f := range frames(dish) {
		if !yield(f) {
			return
		}
	}
}
```

## Your task

`Archive` opens the pass archive with `open`, yields what `ReadPasses` reads from it, and closes it only after the last line. Close it with `defer` instead, right after `open` has given you the file.

> [!TIP]
> When `open` fails there is no file to close: `Archive` yields the error and returns.
