---
slug: same-length-same-time
title: Same Length, Same Time
type: code_tests
position: 4
difficulty: 7
topics: crypto, security, strings
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Thursday. The night shift opens the ops console's maintenance tools with a phrase they type, and the console checks it with `==`. Mirela times it: "It answers sooner the sooner you're wrong. Give me enough tries and I'll read your phrase out piece by piece."

`==` stops at the first difference it finds. `subtle.ConstantTimeCompare` takes the same time for any two slices of one length, but answers at once when the lengths differ, which gives the length away. Hash both phrases first, and every comparison is of 32 bytes:

```go
a := sha256.Sum256([]byte("short"))
b := sha256.Sum256([]byte("a much longer phrase"))
fmt.Println(len(a), len(b)) // 32 32
```

## Your task

`samePhrase` compares with `==`. Make it hash both phrases with `sha256.Sum256` and compare the digests with `subtle.ConstantTimeCompare`, which returns 1 for equal slices. The test reads `main.go`: no `==` or `!=` may compare the phrases.
