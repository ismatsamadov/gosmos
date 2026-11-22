---
slug: leading-zeros
title: Leading Zeros
type: code_tests
position: 3
difficulty: 3
topics: formatting, strings
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:07:33Z
---
The Board shows Brine-3, a customer's satellite, coming into view over Lantern Atoll at `4:7`. Ingrid wants `04:07`: two digits for the hour and two for the minutes, every time.

`fmt.Sprintf` formats exactly as `Printf` does, but returns the text as a string instead of printing it. A `0` before the width pads with zeros instead of spaces:

```go
id := fmt.Sprintf("pipit-1/PEL/%05d", 4182)
fmt.Println(id) // pipit-1/PEL/04182
```

## Your task

Make `clock` in `main.go` return its time as `hh:mm`. `main_test.go` lists the times it checks.
