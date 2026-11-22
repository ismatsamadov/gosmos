---
slug: streams-say-goodbye
title: Streams Say Goodbye
type: code_tests
position: 1
difficulty: 9
topics: channels, http, mutexes
created: 2026-11-07T06:23:24Z
modified: 2026-11-07T06:23:24Z
---
Tonight's checklist gives the Board's deploy at 05:00 two seconds, and Ingrid won't have forty open streams hold it up.

`Shutdown` stops taking connections, then waits for each one to go idle, and a stream never does. `RegisterOnShutdown` names a function for `Shutdown` to run as it starts: closing the hub there ends every stream with a goodbye, and its connection goes idle.

## Your task

1. In `hub.go`, write `Close`: close every channel and delete it from `h.subs`, and mark the hub closed, so that `Subscribe` from then on returns a channel already closed. A second `Close` changes nothing.
2. In `run.go`, make `newServer` register `b.hub.Close` with `RegisterOnShutdown`.

Each screen then reads:

```text
event: bye
data: {}
```

and the end of its stream.
