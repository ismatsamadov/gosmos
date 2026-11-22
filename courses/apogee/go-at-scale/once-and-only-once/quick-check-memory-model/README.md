---
slug: quick-check-memory-model
title: Quick Check: The Memory Model
type: quiz
position: 7
difficulty: 8
topics: goroutines, mutexes
created: 2026-11-11T12:51:43Z
modified: 2026-11-11T12:51:43Z
---
Friday 23 March. All week a question has stood on Nkem's whiteboard; now she wants an answer: "Why does a program without a race behave?"

Under it she has written two ways for two of the Board's screens to count themselves:

```go
var screens int // A
go func() { screens++ }()
screens++

var count atomic.Int64 // B
go func() { count.Add(1) }()
count.Add(1)
```

## Your task

Answer Nkem's three questions: whether snippet A has a data race, why B ends with `count` at 2 once both have run, and what a program with a data race may do.
