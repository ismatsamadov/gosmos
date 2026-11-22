---
slug: close-everything-before-leaving
title: Close Everything Before Leaving
type: code_tests
position: 5
difficulty: 9
topics: goroutines, networking
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Thursday, rehearsing Friday 28 April's deploy, when a new dishd takes the old one's place: told to stop, the old one never does, and keeps both controllers and three watchers.

A clean stop goes in order. Close the listener, so `Accept` returns and nobody new comes in. Close every connection, so every blocked `Read` returns and its goroutine can end. Then wait for those goroutines. Waiting first waits forever: the readers it waits for sit on connections nobody has closed.

```go
go func() {
	_, err := conn.Read(buf) // nothing ever comes
	fmt.Println(err)
	// ...: use of closed network connection
}()
conn.Close() // the Read returns at once
```

## Your task

`Shutdown` closes only the listener. Inside `s.once`, also close every watcher with `s.Watchers.CloseAll()` and every client in `s.Dishes`, then wait for `s.wg`.

> [!NOTE]
> `s.once` runs its function the first time only, so a second `Shutdown` returns at once.
