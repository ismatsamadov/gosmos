---
slug: check-it-on-arrival
title: Check It on Arrival
type: code_tests
position: 2
difficulty: 9
topics: http, security, time
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Yusra is sending customers a sample receiver, and Hollis will copy it line for line into Furrowlight's program.

A receiver believes nothing in a request until it checks out. First the time: more than five minutes from the receiver's own clock, either way, and the request may be an old one sent again. Then the signature, computed again from the bytes that arrived and compared with `hmac.Equal`, which takes as long whichever byte differs. Only then is the body decoded:

```go
sent, _ := strconv.ParseInt("1831732073", 10, 64)
now := time.Unix(1831732134, 0)
fmt.Println(now.Unix() - sent) // 61
```

## Your task

`verifyEvent` decodes whatever arrives. Make it:

1. Read at most 64 KiB of the body, then refuse an `Apogee-Timestamp` that is missing or more than 300 seconds from `now`.
2. Refuse unless `hmac.Equal` finds `Apogee-Signature` equal to `signature(secret, ts, body)`.
3. Only then decode, and return a zero `Event` with every error.
