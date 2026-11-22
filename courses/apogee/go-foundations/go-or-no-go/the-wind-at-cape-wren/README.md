---
slug: the-wind-at-cape-wren
title: The Wind at Cape Wren
type: code_tests
position: 6
difficulty: 4
topics: conditionals, switch
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Nobody lives at Cape Wren, the network's windiest site, so its dish decides for itself: it tracks normally under 50 km/h, tracks with a warning from 50, and stows from 72.

A `switch` with no value after it treats each `case` as a condition. Go tries them from the top and runs the first that is true, like an `else if` chain without the repetition:

```go
switch {
case tempC < -10:
	fmt.Println("too cold")
case tempC > 45:
	fmt.Println("too hot")
default:
	fmt.Println("within limits")
}
```

## Your task

Make `windCall` return `"track"`, `"track with warning"` or `"stow"` for a wind of `windKmh`, using a `switch` with no value.
