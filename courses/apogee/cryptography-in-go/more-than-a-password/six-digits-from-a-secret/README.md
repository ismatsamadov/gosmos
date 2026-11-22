---
slug: six-digits-from-a-secret
title: Six Digits From a Secret
type: code_tests
position: 3
difficulty: 9
topics: bytes, crypto, formatting
created: 2026-11-13T06:18:24Z
modified: 2026-11-13T06:18:24Z
---
Monday 30 October. Mirela's plan for Friday: every engineer's phone gets an authenticator app, which shows a new six-digit code every 30 seconds, made from a secret of 20 random bytes that the console holds too. A password without the phone gets no code.

The recipe is RFC 6238, the published standard for such codes, which it calls TOTP. A code is the HMAC-SHA-1 of a counter, the number of 30-second steps since 1970 as 8 big-endian bytes, cut down by dynamic truncation: the low 4 bits of the HMAC's last byte pick an offset; the 4 bytes from there, top bit cleared, make a number; its last six digits are the code, leading zeros kept.

```go
b := binary.BigEndian.AppendUint64(nil, 1)
fmt.Printf("% x\n", b)
// 00 00 00 00 00 00 00 01
fmt.Printf("%06d\n", 4182) // 004182
```

## Your task

Make `hotp` return the code for `counter`: `hmac.New(sha1.New, secret)` over the counter's 8 bytes, then dynamic truncation, as six digits. `totp`, given, passes it the steps.

> [!NOTE]
> SHA-1 no longer makes safe fingerprints, but inside HMAC it still holds, and authenticator apps use it unless told otherwise.
