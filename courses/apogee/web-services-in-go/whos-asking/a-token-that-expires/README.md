---
slug: a-token-that-expires
title: A Token That Expires
type: code_tests
position: 4
difficulty: 8
topics: security, strings, time
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Mirela wants API keys to cross the network once an hour, not with every request. So a customer's program will trade its key for a token: a short line, signed by the service, that names the customer and when it stops working.

A token is `<customer>.<expiry>.<signature>`. The expiry is in Unix seconds, counted from 1970; the signature is the 64 hex digits of an HMAC-SHA256, under the service's own token key, of `<customer>.<expiry>`. Anyone can read the first two parts, but without the key nobody can change them and sign again:

```go
mac := hmac.New(sha256.New, []byte("lab key"))
mac.Write([]byte("pipit-2|118|SET_MODE|SAFE"))
sig := hex.EncodeToString(mac.Sum(nil))
fmt.Println(len(sig)) // 64
```

## Your task

`signToken` returns `""`. Return the token for `customer` that expires at `expires`, signed with `key`. Furrowlight's token expiring at 11:00 UTC today, 5 January, starts `furrowlight.1830682800.`

> [!WARNING]
> Sign the customer and the expiry together: a signature over the customer alone would let anyone move the expiry and keep a token forever.
