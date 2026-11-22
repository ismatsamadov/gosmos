---
slug: compare-without-a-clue
title: Compare Without a Clue
type: code_tests
position: 1
difficulty: 8
topics: security, strings
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Mirela times the key check from her desk. A check that stops at the first wrong byte answers a little later for every byte a guess gets right, so a patient guesser could find a key one byte at a time.

`subtle.ConstantTimeCompare` compares two byte slices in a time that depends only on their length, never on where they differ. It returns 1 when they are equal and 0 when they aren't:

```go
a := []byte("pipit-2|118|SET_MODE|SAFE")
b := []byte("pipit-2|118|SET_MODE|SAFE")
fmt.Println(subtle.ConstantTimeCompare(a, b))
// 1
fmt.Println(subtle.ConstantTimeCompare(a, b[:7]))
// 0
```

## Your task

For today, `keys` holds each customer's whole API key under its customer ID. `checkKey` returns `""` and false. Take the customer ID before the key's first dot, find that customer's key in `keys`, and compare the whole key with it using `subtle.ConstantTimeCompare`. Return the ID and true only when they match.

> [!NOTE]
> `==` gives the same answers, so the tests also read `auth.go` and check that `checkKey` calls `subtle.ConstantTimeCompare`.
