---
slug: quick-check-bubbles
title: Quick Check: Bubbles
type: quiz
position: 2
difficulty: 8
topics: goroutines, testing, time
created: 2026-11-11T16:51:57Z
modified: 2026-11-11T16:51:57Z
---
Tuesday afternoon, Nkem writes a bubble on the whiteboard and asks what `time.Now()` says in there, and when it moves.

`quit` is a channel made inside the bubble, and nothing ever sends on it:

```go
synctest.Test(t, func(t *testing.T) {
	fmt.Println(time.Now().UTC())
	go func() { <-quit }()
	time.Sleep(10 * time.Minute)
	fmt.Println(time.Now().UTC())
})
```

## Your task

Answer Nkem's three questions: what the first line prints, what a real network read in the goroutine would do to the `Sleep`, and what happens when the function returns with the goroutine still waiting.
