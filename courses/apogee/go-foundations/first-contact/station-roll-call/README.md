---
slug: station-roll-call
title: Station Roll Call
type: code
position: 7
difficulty: 2
topics: printing, statements
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
Friday. Your message came back on Pipit-1's beacon this morning. At the team's short daily meeting, or stand-up, Ingrid wants all five stations on the Board, the operations room's big screen.

A report is a run of `fmt.Println` calls, one per line, in the order the lines should appear. Lines of the same shape get their values in the same order:

```go
fmt.Println("Pass log")
fmt.Println("pipit-1", "AOS", "04:12")
fmt.Println("pipit-1", "LOS", "04:21")
```

## Your task

Add the other four stations below Ystmark, so the program prints:

```text
Station roll call
YST Ystmark up
PEL Port Ellery up
SQT Salar Quieto up
LAN Lantern Atoll up
CWR Cape Wren stowed
```

> [!NOTE]
> A stowed dish is parked safely, usually against wind. Cape Wren's dish stows itself when the wind reaches 72 km/h.
