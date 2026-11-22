---
slug: a-list-of-stations
title: A List of Stations
type: code
position: 0
difficulty: 3
topics: slices
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Nkem reads your roll call code: five station codes in five constants. "And when we open a sixth station?"

A slice holds a list of values of one type. `s[i]` is the element at index `i`, counting from 0, and `len(s)` is how many elements there are:

```go
bands := []string{"UHF", "S", "X"}
fmt.Println(bands[0])   // UHF
fmt.Println(len(bands)) // 3
```

## Your task

Replace the five constants with one slice named `codes`, holding the codes in the same order. Then make the three lines print the first code, the last code and the count in place of each `?`. Find the last code with `len` too, so a sixth station needs no other change.

> [!WARNING]
> The compiler doesn't check slice indexes against the length: `bands[3]` builds, but when that line runs, the program stops with a panic, a run-time error.
