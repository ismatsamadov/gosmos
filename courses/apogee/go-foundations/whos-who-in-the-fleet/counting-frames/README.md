---
slug: counting-frames
title: Counting Frames
type: code
position: 2
difficulty: 4
topics: loops, maps
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:07:33Z
---
Ingrid wants the hourly report to say how many frames of satellite status data, or telemetry, each satellite sent. The feed hands you one satellite ID per frame.

A missing key reads as its zero value, so `m[key]++` works even on a key the map has never seen: the count starts at 0 and becomes 1.

```go
perStation := map[string]int{}
perStation["LAN"]++
perStation["LAN"]++
fmt.Println(perStation) // map[LAN:2]
```

`fmt` prints a map with its keys sorted, so the report reads the same every hour. A `range` over a map makes no such promise: its order can change from run to run.

## Your task

Complete `countBySat` so it counts how many times each ID appears in `ids`. The program should print:

```text
map[brine-3:1 pipit-1:3 thistle:2]
```
