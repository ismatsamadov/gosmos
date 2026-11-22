---
slug: every-change-as-it-happens
title: Every Change as It Happens
type: code_tests
position: 3
difficulty: 8
topics: channels, context, select
created: 2026-11-07T06:23:23Z
modified: 2026-11-07T06:23:23Z
---
The night shift wants a stowed dish on the screen within a second of the refresh that saw it. Nkem gave `Refresh` a channel this morning: while a stream is open, every new view goes on `b.updates`, which has room for 16.

`select` waits on several channel operations at once and runs the first one that is ready:

```go
for {
	select {
	case r := <-readings:
		fmt.Println(r.Station)
	case <-ctx.Done():
		return
	}
}
```

## Your task

In `serveEvents`, after the first event, loop: write each view that arrives on `b.updates` with `writeEvent`, and return when the request's context is done or a write fails. A refresh with nothing new sends nothing, so the screen reads only new versions:

```text
id: 1
id: 2
id: 3
```
