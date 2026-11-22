---
slug: check-the-token
title: Check the Token
type: code_tests
position: 5
difficulty: 9
topics: errors, security, strings, time
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
A Brinewatch engineer asks Mirela what stops them from editing their token to say `cirrine`. Mirela: "The signature. Your check accepts their token; now show me the edited one it refuses."

Trust a token in three steps, in this order: it has three parts; its signature, decoded from hex, equals the HMAC computed again over the first two parts, compared with `hmac.Equal`, which takes the same time wherever they differ; its expiry is still ahead of `now`. Until the signature checks out, nothing in the token is worth believing.

```go
parts := strings.Split("stow.CWR.74", ".")
fmt.Println(len(parts), parts[1]) // 3 CWR
_, err := hex.DecodeString("zz")
fmt.Println(err != nil) // true
```

## Your task

`verifyToken` returns `""` and nil. Return the customer a valid token names. Refuse a token without exactly three parts, or whose signature isn't the one `tokenMAC` computes, with `unauthorized("invalid token")`, and one whose expiry isn't after `now` with `unauthorized("token expired")`.

> [!NOTE]
> At 11:00:00 exactly, the token that expires at 11:00 has expired.
