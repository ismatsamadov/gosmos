---
slug: the-board
title: The Board
type: code
position: 6
difficulty: 4
topics: formatting, printing
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:28Z
---
Friday, 03:40 UTC, the shared time standard Apogee uses at every station. Furrowlight arrives at nine, and the Board still prints the next four passes as a jumble of numbers.

One format string can combine everything this chapter covered, each verb with its own width, decimals or zero padding:

```go
fmt.Printf("%-10s %02d:%02d %5.1f km/h\n",
	"Cape Wren", 9, 5, 51.66)
// Cape Wren  09:05  51.7 km/h
```

## Your task

Replace the five `Println` calls in `main` with `Printf` calls:

1. The header uses `"%-20s %-8s %-6s %6s  %s\n"`.
2. Each pass uses `"%-20s %-8s %02d:%02d  %6.1f  %s\n"`, with its six values after it, in the same order.

The Board then reads:

```text
PASS                 STATION  AOS    MAX EL  BAND
pipit-1/PEL/04182    PEL      04:12    61.4  UHF
furrow-a/SQT/11873   SQT      05:03    23.0  X
thistle/LAN/02218    LAN      06:47     8.9  UHF
cirrine-1/YST/07455  YST      07:30    44.2  S
```
