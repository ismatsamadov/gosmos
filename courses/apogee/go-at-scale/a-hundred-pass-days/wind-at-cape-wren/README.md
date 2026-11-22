---
slug: wind-at-cape-wren
title: Wind at Cape Wren
type: code_tests
position: 6
difficulty: 9
topics: context, goroutines, slices
created: 2026-11-12T03:18:38Z
modified: 2026-11-12T03:18:38Z
---
Wednesday 27 June. Cape Wren's forecast says 74 km/h on Thursday morning, so Ingrid has passplan 2 rehearse it on Thursday's plan: a stow from 09:47 with no end, then its end at 13:23.

An event's stow with a zero `To` begins at its station; one with a `To` ends the station's stow in progress. The `Replanner` plans every run with every stow it has heard of, so an ended stow still keeps its passes out. Runs in flight read the list while `Run` takes the next event, so a change makes a new slice and leaves the old one alone:

```go
next := slices.Clone(gusts)
next[0].Kmh = 74 // gusts is as it was
gusts = next
```

## Your task

`Run` plans with `r.stows`, which nothing changes, and `heard` returns its list as it was. Make `heard` return a new list: `st` added or, for a stow with a `To`, its station's open stow ended. Then have `Run` pass each event's stow through `heard` into `r.stows` before the run starts.

> [!NOTE]
> The tests rehearse Thursday's real plan: the stow takes Cape Wren's 34 contacts after 09:47 away, and its end leaves Cape Wren 71.
