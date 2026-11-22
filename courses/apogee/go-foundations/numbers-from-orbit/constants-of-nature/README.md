---
slug: constants-of-nature
title: Constants of Nature
type: code
position: 6
difficulty: 3
topics: constants
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Tomás finds Earth's radius typed as 6371, 6378 and 6370 km in three tools. "The Earth isn't going to change," he says, "so pick one and name it once."

`const` names a value that is fixed when the program is compiled: it comes from your code, never from a clock or a file, and nothing can change it later. A constant declared without a type is untyped, so a whole number like `maskDeg` works with `int`s and `float64`s alike:

```go
const maskDeg = 5

elevDeg := 12.5
fmt.Println(elevDeg - maskDeg) // 7.5
```

## Your task

Declare `const earthRadiusKm = 6371` above `main` and use it in place of all three radii. An orbit's radius is Earth's radius plus the altitude.
