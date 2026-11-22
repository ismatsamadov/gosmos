---
slug: columns-that-line-up
title: Columns That Line Up
type: code
position: 2
difficulty: 3
topics: formatting
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:19:30Z
---
On the Board, the column for each station's lowest usable angle above the horizon, its mask, wanders because station names differ in length.

A number between `%` and the verb sets a minimum width. The padding goes on the left, which lines numbers up on the right; with a `-`, it goes on the right, which suits text. Width and decimals combine:

```go
fmt.Printf("[%6s]\n", "PEL")  // [   PEL]
fmt.Printf("[%-6s]\n", "PEL") // [PEL   ]
fmt.Printf("[%6.2f]\n", 7.9)  // [  7.90]
```

## Your task

In `main`, give the station names a width of 13, padded on the right, and the masks a width of 5 with one decimal. Pad the header the same way, so the Board reads:

```text
STATION        MASK
Ystmark         5.0
Salar Quieto    3.0
Lantern Atoll   8.0
```
