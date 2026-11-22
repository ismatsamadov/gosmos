---
slug: counted-not-guessed
title: Counted, Not Guessed
type: code_tests
position: 1
difficulty: 9
topics: goroutines, json, metrics
created: 2026-11-12T09:13:51Z
modified: 2026-11-12T11:24:16Z
---
Tuesday 17 July. Yusra: "TL;DR: the investors will ask how many passes flew, and 'about' is not a number."

A counter only goes up, and many goroutines add to it at once; `atomic.Uint64` from `sync/atomic` does that without a lock:

```go
var heard atomic.Uint64
heard.Add(1) // safe from any goroutine
fmt.Println(heard.Load())
```

## Your task

Count the week as it runs. `reports/metrics.json`, saved after it, holds one JSON object, a count by two labels written as `{"brine-3":{"flown":23}}`:

```text
passes_total      satellite, flown/failed/skipped
failures_total    satellite, reason
frames_total      TLM frames with a good CRC
bad_frames_total  frames with a bad CRC
bad_lines_total   TLM lines the parser refused
commands_total    satellite, accepted/refused
files_landed_total
deliveries_total  customer, delivered/undelivered
```

`ack_wait_seconds` is a histogram of the wait from a command sent to its answer: `bounds` `[0.5,1,2,5,10,30,60,120,300]`; `counts`, one per bound, each wait in the first bound at least as long, then one for longer waits; `count`; `sum_ms`.

The tests compare every number with the simulated world's own counts, for the week and for Thursday at three stations. A 2 s answer counts under 2, Ystmark's 2.7 s ones over the backhaul under 5.
