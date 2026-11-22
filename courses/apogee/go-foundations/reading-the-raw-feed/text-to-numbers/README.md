---
slug: text-to-numbers
title: Text to Numbers
type: code_tests
position: 2
difficulty: 5
topics: conversions, errors
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
One of Thistle's frames read `batt_v=7.9l`, with a letter l where the 1 belongs, and the parser stored 0 V without a word.

`strconv.ParseFloat(s, 64)` turns text into a `float64`; the 64 asks for that size. When the text isn't a number, it returns 0 and an error, and that error is the only sign that anything went wrong. Below, the last character of `"437.25O"` is the letter O, typed for a zero:

```go
mhz, err := strconv.ParseFloat("437.25O", 64)
fmt.Println(mhz) // 0
fmt.Println(err)
// strconv.ParseFloat: parsing "437.25O":
//     invalid syntax
```

## Your task

In `battery`, stop throwing away the error from `strconv.ParseFloat`. When there is one, return 0 and the error, wrapped with `%w` in a message that names `batt_v`.
