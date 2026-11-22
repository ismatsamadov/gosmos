---
slug: sign-what-you-send
title: Sign What You Send
type: code_tests
position: 1
difficulty: 8
topics: http, security
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Mirela reads yesterday's webhook and leaves one line: "Anyone can POST to Furrowlight's address. Only we can sign."

So every request carries proof: an HMAC-SHA256, under the customer's webhook secret, of the request's Unix time, a dot, then the exact bytes of its body. The time goes in `Apogee-Timestamp` and the signature, as 64 hex digits, in `Apogee-Signature`. A receiver that knows the secret computes the same digits; anyone else can only guess them. You write to an HMAC as to a file:

```go
mac := hmac.New(sha256.New, key)
mac.Write([]byte("pipit-2|118|"))
mac.Write([]byte("PING"))
fmt.Println(hex.EncodeToString(mac.Sum(nil)))
```

## Your task

1. `signature` returns `""`. Return the hex HMAC-SHA256, under `secret`, of `ts` in decimal, a dot, then `body`.
2. In `post`, set `Apogee-Timestamp` to `ts` and `Apogee-Signature` to the signature of the bytes you send.

> [!WARNING]
> Sign the time with the body. A signature of the body alone lets anyone who copied an old request send it again under a new time.
