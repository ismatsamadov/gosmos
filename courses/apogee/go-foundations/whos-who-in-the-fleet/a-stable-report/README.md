---
slug: a-stable-report
title: A Report in Stable Order
type: code_tests
position: 5
difficulty: 5
topics: maps, sorting
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:28Z
---
Ingrid notices that the hourly report lists the satellites in a different order every hour. It ranges over a map, and a map keeps no order of its own.

To visit a map's keys in order, sort them first. `maps.Keys(m)` hands out the keys one at a time, and `slices.Sorted` collects them into a sorted slice:

```go
masks := map[string]float64{
	"SQT": 3, "CWR": 5, "LAN": 8,
}
codes := slices.Sorted(maps.Keys(masks))
fmt.Println(codes) // [CWR LAN SQT]
```

Add `"maps"` and `"slices"` to the imports to use them.

## Your task

Complete `report` so it returns one `"id: count"` line per satellite, such as `"pipit-1: 3"`, sorted by ID. For an empty map, return no lines.
