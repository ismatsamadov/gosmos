---
slug: thistles-last-pass
title: Thistle's Last Pass
type: code
position: 9
difficulty: 8
topics: channels, goroutines, maps, sorting
created: 2026-11-01T07:36:47Z
modified: 2026-11-06T16:51:26Z
---
Thistle re-enters tonight: it returns into Earth's atmosphere at the end of its life. On its last pass four stations listen at once, and Professor Marr wants to know which beacons were heard, and by whom.

Concurrent results arrive in whatever order the goroutines happen to run. Collect them first, here in a map, then sort them before printing, so the report reads the same on every run.

```go
// seen lists the stations in arrival order.
seen := []string{"SQT", "LAN", "PEL"}
slices.Sort(seen)
fmt.Println(strings.Join(seen, " "))
// LAN PEL SQT
```

## Your task

1. `merge` passes on Port Ellery's feed only. Make it listen to every feed, one goroutine each, and close its output once all are done.
2. Run the program a few times: the station lists change order. Sort each list before printing it.

> [!NOTE]
> Thistle never rises over Ystmark, so YST isn't listening.
