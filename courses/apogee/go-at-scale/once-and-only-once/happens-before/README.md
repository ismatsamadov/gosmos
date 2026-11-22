---
slug: happens-before
title: Happens Before
type: code
position: 4
difficulty: 8
topics: goroutines, mutexes
created: 2026-11-11T12:51:42Z
modified: 2026-11-11T12:51:42Z
---
Wednesday 21 March. The test Board, a screen beside Ingrid's that shows passplan 2's plans, read a plan the moment its ready flag went up and showed version 0 with no contacts: the planner had raised the flag first and filled the plan after.

Go's memory model says when one goroutine sees another's writes. An atomic store happens before every load that sees it, and so does everything its goroutine did before the store, but not what it does after. Here the planner waits until the Board has looked before filling the plan, so the wrong order shows on every run.

```go
go func() {
	wind = 64          // before the Store: seen
	posted.Store(true)
}()
for !posted.Load() { // waits for the Store
}
fmt.Println(wind) // 64, on every run
```

## Your task

Move the two lines that fill the plan above `ready.Store(true)`, so the test Board prints the plan:

```text
plan v7: 2016 contacts
```

> [!NOTE]
> Neither version has a data race: `go run -race .` in a terminal reports nothing. Only the order was wrong.
