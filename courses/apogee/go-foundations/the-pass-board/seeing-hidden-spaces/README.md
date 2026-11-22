---
slug: seeing-hidden-spaces
title: Seeing Hidden Spaces
type: code
position: 5
difficulty: 3
topics: formatting, types
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Every lookup for Salar Quieto fails, yet its code looks right on the Board, and its mask prints as `3`, like a whole number. Paz suspects a typo.

`%s` and `%v` print what a value looks like; two other verbs show what it is. `%q` puts a string in double quotes, so spaces at its ends show, and writes a character such as a tab as an escape like `\t`; `%T` prints a value's type:

```go
band := "S\t"
fmt.Printf("%q\n", band) // "S\t"
fmt.Printf("%T\n", band) // string
fmt.Printf("%T\n", 437)  // int
```

## Your task

1. Print the three station codes with `%q`, so the stray space shows.
2. Make the last line print the mask's type instead of its value: `mask type: float64`.

> [!TIP]
> When two strings look equal but don't match, print both with `%q`.
