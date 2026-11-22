---
slug: the-last-seconds-kept
title: The Last Seconds, Kept
type: code_tests
position: 6
difficulty: 10
topics: errors, files, tracing
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:39Z
---
Tuesday 20 February. Nkem asks whether the stations could run a trace all night. Ingrid: "Negative. Five stations, every night: nobody reads that. A failed pass leaves its last seconds behind."

A flight recorder runs a trace in memory and holds only the last stretch of it: about the last `MinAge` of events, unless that would take more than about `MaxBytes`. Nothing reaches the disk until something asks: `WriteTo` writes what it holds as an ordinary trace, and the recorder goes on recording:

```go
cfg := trace.FlightRecorderConfig{
	MinAge:   time.Minute,
	MaxBytes: 16 << 20,
}
fr := trace.NewFlightRecorder(cfg)
err := fr.Start()
fmt.Println(err, fr.Enabled()) // <nil> true
```

## Your task

`watchPass` only calls `run`. When `run` fails, write the recorder's snapshot with `fr.WriteTo` into a new file at the given path, and return `run`'s error joined with any error from saving it (`errors.Join`). When `run` succeeds, write nothing and return nil.

> [!TIP]
> In a terminal, `go tool trace flight.trace` opens the last seconds before the failure.
