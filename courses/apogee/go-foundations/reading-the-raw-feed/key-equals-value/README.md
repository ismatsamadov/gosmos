---
slug: key-equals-value
title: Key Equals Value
type: code_tests
position: 1
difficulty: 5
topics: multiple-returns, strings
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:27Z
---
The pager wants single readings out of Thistle's lines: from `batt_v=7.91`, the key `batt_v` and the value `7.91`.

`strings.Cut(s, sep)` splits `s` at the first `sep` only. It returns the text before, the text after, and whether `sep` was there at all:

```go
id := "pipit-1/PEL/04182"
before, after, found := strings.Cut(id, "/")
fmt.Println(before, after, found)
// pipit-1 PEL/04182 true

_, _, found = strings.Cut("thistle", "/")
fmt.Println(found) // false
```

## Your task

Inside the loop in `field`, cut each field at its first `=`. When the field has an `=` and the part before it is `key`, return the part after it and `true`.
