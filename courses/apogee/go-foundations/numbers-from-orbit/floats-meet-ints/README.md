---
slug: floats-meet-ints
title: Floats Meet Ints
type: code_tests
position: 5
difficulty: 3
topics: conversions, types
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:07:33Z
---
Yusra forwards a short summary from Furrowlight, a customer that photographs farmland: their dashboard shows 333 MB per pass, but their own count says 333.33. She needs the two numbers to agree.

Dividing an `int` by an `int` gives an `int`: the fraction is gone at once, and converting the result later can't bring it back. `float64(x)` converts `x` to a `float64`, and Go only does arithmetic on two values of the same type:

```go
passes := 4
fmt.Println(30 / passes)          // 7
fmt.Println(30 / float64(passes)) // 7.5
```

A bare number like `30` takes whichever type it meets; more on that soon.

## Your task

`averageMB` divides first and converts afterwards. Change its `return` line so it converts before dividing, and the tests in `main_test.go` pass.
