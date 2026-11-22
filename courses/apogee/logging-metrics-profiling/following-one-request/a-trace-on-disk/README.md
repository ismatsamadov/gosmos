---
slug: a-trace-on-disk
title: A Trace on Disk
type: code_tests
position: 0
difficulty: 9
topics: files, tracing
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:39Z
---
Tuesday 13 February, the lab bench. The replay begins as the night did, three of the old uplinkd's four workers already held. Nkem: "This time we watch it happen."

An execution trace, from `runtime/trace`, shows it all: when each goroutine starts, waits, wakes and ends, and what it waited for, to the microsecond. `trace.Start(w)` starts writing one to `w`, and `trace.Stop()` ends it once everything is written. Every trace begins the same way:

```go
var buf bytes.Buffer
if err := trace.Start(&buf); err != nil {
	log.Fatal(err)
}
trace.Stop()
fmt.Printf("%q\n", buf.Bytes()[:5]) // "go 1."
```

## Your task

`traceTo` only runs `work`. Make it create the file at the path it's given, start a trace into it, run `work`, stop the trace and close the file, returning the first error.

> [!TIP]
> In a terminal, `go tool trace trace.out` opens the trace in your browser, every goroutine of the replay on one timeline.
