---
slug: above-the-mask
title: Above the Mask
type: code_tests
position: 1
difficulty: 3
topics: conditionals
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
The scheduler books every predicted pass, even ones that never clear the station's mask. Ingrid wants each pass called first, with the numbers behind the call: GO if its highest point is at or above the mask, NO-GO if not.

`else` gives an `if` a second block, which runs only when the condition is false. Exactly one of the two blocks runs, so `if` and `else` are a natural way to give a variable one of two values:

```go
var payload string
if tempC > 45 {
	payload = "off"
} else {
	payload = "on"
}
fmt.Println("payload:", payload)
```

## Your task

In `elevationCall`, replace the TODO with an `if` and an `else` that set `call` to `"GO"` when `maxElevDeg` is at or above `maskDeg`, and to `"NO-GO"` otherwise. The `Sprintf` below them already adds the numbers.
