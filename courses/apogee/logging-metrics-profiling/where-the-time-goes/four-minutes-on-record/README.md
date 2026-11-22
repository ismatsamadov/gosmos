---
slug: four-minutes-on-record
title: Four Minutes on Record
type: code_tests
position: 0
difficulty: 8
topics: files, profiling
created: 2026-11-10T15:36:39Z
modified: 2026-11-11T17:33:38Z
---
Tuesday morning, at the backend desks. Nobody knows yet where passplan, the program that predicts every pass, spent its four minutes, and Nkem won't guess: "Record it. Then we'll know what to fix."

A profile records where a program spends its time. Between `pprof.StartCPUProfile` and `pprof.StopCPUProfile`, from `runtime/pprof`, Go notes 100 times a second which functions are running, and the stop writes it all to the file the start was given. Only one profile runs at a time, so a start can fail:

```go
pprof.StartCPUProfile(io.Discard)
err := pprof.StartCPUProfile(io.Discard)
fmt.Println(err)
pprof.StopCPUProfile()
// cpu profiling already in use
```

## Your task

Make `profileTo` create the file at the path it's given, start a CPU profile into it, run `work`, stop the profile and close the file. Return the first error, and leave `work` unrun when the file or the profile can't start. The test profiles passplan's plan of Tuesday into `cpu.prof`.

> [!TIP]
> `go tool pprof -top cpu.prof`, run in a terminal in the lesson folder, lists the functions that took the time.
