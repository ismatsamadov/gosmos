---
slug: every-screen-gets-every-change
title: Every Screen Gets Every Change
type: code_tests
position: 0
difficulty: 8
topics: channels, maps, mutexes
created: 2026-11-07T06:23:23Z
modified: 2026-11-07T06:23:23Z
---
Wednesday. Overnight the meeting room got the 03:58 change and the kitchen didn't: two streams were taking turns on one channel. Nkem moved the stream, `Refresh` and `/health` onto a hub this morning, and left its methods for you.

A value sent on a channel reaches one receiver. To reach every screen, give each stream a channel of its own and send every view on all of them. Streams join and leave from many goroutines at once, so the hub keeps the channels in a map behind a mutex, keyed by the end each stream receives from:

```go
subs := map[<-chan int]chan int{}
ch := make(chan int, 4)
subs[ch] = ch
```

## Your task

In `hub.go`, write:

1. `Subscribe`: a new channel with room for 4 views, kept in `h.subs`.
2. `Unsubscribe`: if `h.subs` holds the channel, delete it and close it.
3. `Publish`: send the view on every channel.
4. `Len`: how many channels `h.subs` holds.

> [!NOTE]
> Until the hub works, the earlier stream tests fail too.
