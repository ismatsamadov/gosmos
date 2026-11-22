---
slug: quick-check-structs
title: Quick Check: Structs
type: quiz
position: 4
difficulty: 3
topics: structs, variables
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Your types go into review tomorrow. Nkem sends three questions first, one at a time, and the first two are about this code:

```go
type Station struct {
	Code    string
	MaskDeg float64
}
var s Station
a := Station{Code: "PEL", MaskDeg: 5}
b := a
b.MaskDeg = 9
```

Read the code one line at a time, and answer from what it does, not from memory.
