---
slug: when-the-loop-stops
title: When the Loop Stops
type: code_tests
position: 1
difficulty: 8
topics: functions, iterators
created: 2026-11-11T21:35:30Z
modified: 2026-11-11T21:35:30Z
---
Tuesday 27 March. Hollis asked the booking service, bookingd, for Furrow-A's first three passes of July over Salar Quieto (SQT). The loop serving him stopped at three; the draft iterator behind bookingd's catalog didn't, and bookingd panicked.

When a loop over an iterator stops early, with `break` or `return`, the call to `yield` that ran its body returns false. From then on the iterator must not call `yield` again, or Go panics. Inside two loops, return from the function: a `break` would leave only the inner one.

```go
for _, feed := range feeds {
	for _, f := range feed.frames {
		if !yield(f) {
			return // both loops end here
		}
	}
}
```

## Your task

`Passes` asks `predict` for one day at a time and yields each pass, ignoring what `yield` returns. Return as soon as it returns false, before another day is predicted.

> [!NOTE]
> Go's panic reads "runtime error: range function continued iteration after function for loop body returned false".
