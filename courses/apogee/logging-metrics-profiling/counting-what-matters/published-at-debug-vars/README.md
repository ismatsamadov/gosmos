---
slug: published-at-debug-vars
title: Published at Debug Vars
type: code_tests
position: 5
difficulty: 8
topics: http, json, metrics
created: 2026-11-10T15:36:37Z
modified: 2026-11-10T15:36:37Z
---
Tuesday 12 December, 09:00, an hour before the counters go to every station. The pager will read each program's numbers every five minutes, and the Board reads them too, all at one address on the program's debug port, the private port where it answers the team: `GET /debug/vars`.

The `expvar` package serves every variable published with `expvar.Publish` as one JSON object at `/debug/vars`, beside two of its own, `cmdline` and `memstats`. An `expvar.Func` is a function whose result is turned into JSON each time someone asks, so a reading is never stale:

```go
expvar.Publish("site", expvar.Func(func() any {
	return "LAN"
}))
// GET /debug/vars, in part:
// "site": "LAN"
```

## Your task

`publish` publishes nothing. Inside `publishOnce.Do`, publish three `expvar.Func` values: `passes_total`, returning `m.Passes.Snapshot()`; `ack_wait_seconds`, returning `m.AckWait.Snapshot()`; and `goroutines`, returning `runtime.NumGoroutine()`.

> [!NOTE]
> `expvar.Publish` panics on a name it already has, so `publishOnce` keeps a second call to `publish` harmless.
