---
slug: a-hundred-days-under-load
title: A Hundred Days Under Load
type: code_tests
position: 7
difficulty: 10
topics: context, functions, iterators
created: 2026-11-12T03:18:38Z
modified: 2026-11-12T10:15:10Z
---
Wednesday 27 June, 21:00: Furrowlight's twenty-four are being heard for the first time, and on its last evening in shadow, passplan 2 plans a hundred days for all eighty satellites. When Nkem cancels a run partway, it goes on predicting.

`planStation` checks its context every 1,000 candidates, but a stream can predict many days before it yields one. Brinewatch's satellites never climb ten degrees over the three stations in the far north, Ystmark, Tirrin Rock and Ilvane Spit, so there the merge's first pull asks the pass model for every day of the window, satellite after satellite. The check belongs inside the streams, before each day's prediction. A function of your own can stand in front of the model:

```go
quiet := func(city string) Forecast {
	if ctx.Err() != nil {
		return Forecast{} // canceled: ask no more
	}
	return fetch(city)
}
```

## Your task

In `planStation`, hand `stationPasses` a model of your own instead of `passModel`: it returns no passes once `ctx` has ended, and asks `passModel` otherwise.

> [!NOTE]
> Nkem's test plans the first week from the hand-over, 14,282 contacts, then cancels a run when Monday 2 July is first asked for.
