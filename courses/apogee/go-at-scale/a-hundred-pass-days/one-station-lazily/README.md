---
slug: one-station-lazily
title: One Station, Lazily
type: code_tests
position: 0
difficulty: 9
topics: context, iterators
created: 2026-11-12T03:18:38Z
modified: 2026-11-12T10:15:10Z
---
Nkem starts you with one station's plan: every planning run makes twenty of them, one for each station.

Every piece is there already: `stationPasses` merges a station's streams into candidates a day at a time, `stowedDuring` says whether a stow covers a pass, and `assign` gives dishes. A sequence of your own joins them. A run can be canceled at any moment, so the sequence watches the context as it goes, at its first candidate and every 1,000th after, as this loop does:

```go
for i, line := range log {
	if i%1000 == 0 && ctx.Err() != nil {
		return ctx.Err()
	}
	send(line)
}
```

## Your task

`planStation` returns nothing. Give `assign` a sequence that ranges over `stationPasses`, asking `passModel`, and yields each candidate no stow covers until `ctx` ends. Return the contacts and the dropped count, or `ctx.Err()` once the context has ended.

> [!NOTE]
> `passModel` is the pass model itself; the tests put their own in its place to cancel a week's run on its second day.
