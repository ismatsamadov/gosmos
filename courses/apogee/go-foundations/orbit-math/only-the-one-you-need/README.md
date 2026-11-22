---
slug: only-the-one-you-need
title: Only the One You Need
type: code
position: 3
difficulty: 3
topics: compiler-errors, multiple-returns
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
The pass banner needs only the minutes of Pipit-1's AOS, and now it won't build: `declared and not used: h`.

Go refuses to build a program in which a function's body declares a variable and never uses it. When you don't need one of a function's results, receive it with the blank identifier `_`, which takes the value and throws it away:

```go
aos, _ := window(252, 9)
fmt.Println("AOS:", aos) // AOS: 252
```

## Your task

Fix the line the compiler names, so the program builds and prints `Minutes past the hour: 12`.
