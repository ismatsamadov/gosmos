---
slug: a-reset-file-worth-stealing
title: A Reset File Worth Stealing
type: code_tests
position: 0
difficulty: 9
topics: crypto, maps, security
created: 2026-11-13T06:18:23Z
modified: 2026-11-13T06:18:23Z
---
Wednesday 25 October. Before any new lock, Mirela checks the console's back door: the reset link it emails to anyone who forgets a password. Each pending link's token sits in the reset file exactly as it was sent. "In July a laptop left with a copy of the user file. Assume the next one holds this file."

Keep a reset token the way the console keeps session tokens: only its hex SHA-256, which `tokenDigest` returns. A copy of the file then holds digests, and a digest pasted into a link is one more wrong token. To check a link, hash its token and compare the two digests in constant time:

```go
a := []byte(tokenDigest("fishbowl"))
b := []byte(tokenDigest("fishbowl!"))
fmt.Println(subtle.ConstantTimeCompare(a, b)) // 0
```

## Your task

1. Make `Start` keep `tokenDigest(token)` as the reset's `Hash`, never the token.
2. Make `Redeem` compare `tokenDigest(token)` with the kept `Hash` using `subtle.ConstantTimeCompare`, with no `==` or `!=` between them.
