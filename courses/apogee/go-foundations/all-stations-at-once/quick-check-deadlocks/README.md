---
slug: quick-check-deadlocks
title: Quick Check: Deadlocks
type: quiz
position: 8
difficulty: 6
topics: channels, goroutines
created: 2026-11-01T07:36:47Z
modified: 2026-11-04T17:07:34Z
---
Ingrid's incident review, the team's account of what went wrong and what to change, lists two snippets that crashed the archiver last night. The first is a whole `main`:

```go
func main() {
	acks := make(chan string)
	acks <- "ACK 118"
	fmt.Println(<-acks)
}
```

The second ran at the end of a pass, where `frames` is a `chan uint32`:

```go
close(frames)
close(frames)
```

The last question asks what a receive from `frames` gives once it is closed and empty.
