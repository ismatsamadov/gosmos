---
slug: a-chain-for-ingrid
title: A Chain for Ingrid
type: code
position: 4
difficulty: 8
topics: iterators, slices, strings
created: 2026-11-11T22:52:57Z
modified: 2026-11-11T22:52:57Z
---
Tuesday 3 April. Hollis asked, in capitals, whether the new stations heard Furrowlight's satellites on their first day, and Ingrid answers before her 09:00 checklist with the first five contacts, as the Board shows them.

Adapters chain: each takes the sequence of the one inside it, so a chain reads from the inside out. Nothing runs until something ranges over the outermost or collects it, and you collect where a slice is needed, as `strings.Join` needs one:

```go
kmh := map[string]int{
	"CWR": 81, "OST": 77, "YST": 40,
}
stows := func(code string) bool {
	return kmh[code] >= 72
}
windy := Filter(maps.Keys(kmh), stows)
fmt.Println(slices.Sorted(windy)) // [CWR OST]
```

## Your task

`main` prints every contact of Monday's schedule as a Board line. Chain `Filter`, `Map` and `Take` (in `adapt.go`) with `furrowAtNew` and `line` instead, then collect and join the lines, to print only the first five Furrowlight contacts at the stations of the first wave:

```text
02:00:24 furrow-a/OST/23203 dish 1
02:47:33 furrow-b/OST/23201 dish 1
03:35:27 furrow-a/OST/23204 dish 1
04:26:40 furrow-a/VSK/23205 dish 1
05:00:58 furrow-b/CAL/23203 dish 1
```
