---
slug: twelve-workers-one-queue
title: Twelve Workers, One Queue
type: code_tests
position: 1
difficulty: 8
topics: channels, goroutines
created: 2026-11-11T12:51:42Z
modified: 2026-11-11T12:51:42Z
---
Tuesday 13 March. Each prediction needs working memory for a day of passes, and Nkem would rather set that up twelve times than eighty: twelve workers that live for the whole run, each taking the next satellite from one queue.

A worker pool is a fixed number of goroutines ranging over one channel of jobs. A worker takes the next job as soon as it has finished the last, so twelve workers hold twelve jobs at most, however many wait in the channel. When the channel is closed and empty, every worker's loop ends.

```go
for range 3 { // three checkers, one queue
	wg.Go(func() {
		for f := range frames {
			check(f)
		}
	})
}
wg.Wait() // frames closed and drained
```

## Your task

`runPool` starts a goroutine for every job. Start exactly `workers` goroutines instead, each ranging over `jobs` and calling `do` for every job it takes, and keep `fail` for the errors.

> [!NOTE]
> The caller closes `jobs` once every satellite is in it; until then, an idle worker waits for more.
