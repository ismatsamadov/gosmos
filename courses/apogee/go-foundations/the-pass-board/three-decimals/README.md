---
slug: three-decimals
title: Three Decimals, Always
type: code
position: 1
difficulty: 2
topics: formatting
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Operators read frequencies to three decimals, but the Board shows `437.25` for Pipit-1 and `2245` for Pipit-2. Ingrid wants every frequency in the same shape.

`%v` prints a value in its default format, which for a float is as short as possible. `%.3f` always prints three decimals, adding zeros or rounding; `%.1f` prints one:

```go
battV := 7.96
fmt.Printf("%v V\n", battV)   // 7.96 V
fmt.Printf("%.3f V\n", battV) // 7.960 V
fmt.Printf("%.1f V\n", battV) // 8.0 V
```

## Your task

Change the verbs in `main` so both frequencies print with three decimals, and the maximum elevation, the pass's highest point above the horizon, with one.
