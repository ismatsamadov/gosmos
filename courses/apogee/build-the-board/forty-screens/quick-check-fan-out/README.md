---
slug: quick-check-fan-out
title: Quick Check: Fan-Out
type: quiz
position: 3
difficulty: 8
topics: channels, select
created: 2026-11-07T06:23:23Z
modified: 2026-11-07T06:23:23Z
---
After the forty-screen rehearsal, Nkem has three questions about channels that many goroutines share:

```go
ch := make(chan View, 4)
go stream(ch) // the kitchen
go stream(ch) // the meeting room
select {
case ch <- v:
default:
}
```

## Your task

Answer three questions: which stream gets a view sent on a shared channel, what `select` with `default` does when the channel is full, and why only the hub closes a stream's channel.
