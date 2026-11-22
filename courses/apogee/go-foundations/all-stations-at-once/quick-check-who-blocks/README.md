---
slug: quick-check-who-blocks
title: Quick Check: Who Blocks?
type: quiz
position: 2
difficulty: 6
topics: channels, goroutines
created: 2026-11-01T07:36:47Z
modified: 2026-11-22T10:08:26Z
---
The Ystmark archiver stalled for ten minutes last night. Nkem wants to be sure you can tell which line of a program waits, and draws two channels on the whiteboard:

```go
a := make(chan string)    // no buffer
b := make(chan string, 2) // room for two values
```

In the first two questions, no goroutine is receiving from either channel yet. The third is about the program as a whole.

## Your task

Answer Nkem's three questions.
