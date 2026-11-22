---
slug: checking-a-login
title: Checking a Login
type: code_tests
position: 4
difficulty: 8
topics: crypto, security
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Wednesday. Mirela made her own record first, and wants to log in to a test copy of the console with the password behind it: "First it lets me in. Then show me who it keeps out."

The console never stores a password, so checking one means deriving again from what was typed, with the record's own count and salt, and comparing the result with the stored key. Compare with `subtle.ConstantTimeCompare` from `crypto/subtle`: it looks at every byte however early two keys differ, so the time a refusal takes says nothing about the stored key. It returns 1 for equal slices:

```go
a := []byte("SET_MODE|SAFE")
b := []byte("SET_MODE|SAVE")
fmt.Println(subtle.ConstantTimeCompare(a, a)) // 1
fmt.Println(subtle.ConstantTimeCompare(a, b)) // 0
```

## Your task

Make `checkPassword` derive a key from `pw` with `pbkdf2.Key`, using `r.Iterations` and `r.Salt`, and report whether it equals `r.Key`, compared with `subtle.ConstantTimeCompare`. Return any error from `pbkdf2.Key`.

> [!NOTE]
> Use the record's count, not 600,000: the lab bench's test account, `bench`, has a record made with 1,000 iterations, and it must keep working until it is replaced.
