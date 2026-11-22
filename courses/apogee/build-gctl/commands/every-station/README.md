---
slug: every-station
title: Every Station
type: code_tests
position: 2
difficulty: 5
topics: slices, strings
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
Before the first pass, Ingrid checks which dishes can move. In high wind a dish parks itself pointing straight up, which ops calls a stow, and it can't follow a satellite until the wind drops.

A slice keeps its order: ranging over `s.stations` visits the stations in the table's order every time, which a map never promises.

## Your task

1. Write `cmdStations`: a line per station, in the table's order, with its code and name, then `stowed:` and the reason for a stowed station, or `stowed` alone when it has no reason.
2. Make any argument an error.
3. Add `stations` to the table.

With Cape Wren stowed, `stations`, `stations PEL` and `help` answer with these lines, among others:

```text
LAN Lantern Atoll
CWR Cape Wren stowed: wind 74 km/h
error: stations takes no arguments
stations  list the stations and which are stowed
```
