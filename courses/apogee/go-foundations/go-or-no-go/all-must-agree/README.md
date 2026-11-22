---
slug: all-must-agree
title: All Must Agree
type: code
position: 3
difficulty: 3
topics: conditionals
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:07:33Z
---
Ingrid's poll starts with one question: can the station follow the satellite across the sky, or track it? Only if the station is up, its dish isn't stowed and the wind is under 72 km/h.

`&&` (and) is true when both sides are, `||` (or) when at least one is, and `!` (not) flips a `bool`. Go reads left to right and stops once the answer is certain: in `a && b`, a false `a` means `b` never runs.

```go
mode, battV := "SAFE", 7.9
needsCare := mode == "SAFE" || battV < 6.6
fmt.Println(needsCare, !needsCare) // true false
```

## Your task

Change the `return` in `canTrack` so it combines all three conditions. Port Ellery's dish is stowed for maintenance and Cape Wren's wind has reached 72 km/h, so only Ystmark should print `true`.
