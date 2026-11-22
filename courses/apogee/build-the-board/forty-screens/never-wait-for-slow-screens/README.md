---
slug: never-wait-for-slow-screens
title: Never Wait for Slow Screens
type: code_tests
position: 1
difficulty: 8
topics: channels, select
created: 2026-11-07T06:23:23Z
modified: 2026-11-07T06:23:23Z
---
At 12:40 Ingrid finds every screen in the building frozen: the reception laptop went to sleep with its stream open, and `Refresh` is waiting for it.

A send on a full channel waits until someone receives. A `select` with a `default` case never waits: when no other case is ready, `default` runs.

```go
select {
case alerts <- a:
default:
	fmt.Println("no room: alert dropped")
}
```

A stream whose 4 slots are full has fallen behind by four views. Drop it: the screen will reconnect and read the Board as it is then.

## Your task

In `Publish`, send each view with `select` and `default`. A channel with no room is deleted from `h.subs` and closed, and its stream ends with:

```text
event: bye
data: {}
```
