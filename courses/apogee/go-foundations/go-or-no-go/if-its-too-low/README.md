---
slug: if-its-too-low
title: If It's Too Low
type: code
position: 0
difficulty: 2
topics: conditionals
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Ingrid reads the pass log: "Three passes, three warnings, one low pass." The warning belongs only under a pass that peaks below 10 degrees.

`if` runs the statements in its braces only when its condition is true. A condition is a `bool`, often from comparing two values with `<`, `<=`, `>`, `>=`, `==` (equal) or `!=` (not equal). Parentheses around it are optional and usually left off, but the braces are required, even around one statement:

```go
tempC := -12.5
if tempC < -10 {
	fmt.Println("battery too cold: heaters on")
}
```

## Your task

In `main`, wrap each of the three warning lines in `if maxElevDeg < 10 { ... }`, so only the low pass keeps its warning.
