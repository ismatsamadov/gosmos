---
slug: signing-the-uplink
title: Signing the Uplink
type: code_tests
position: 9
difficulty: 7
topics: bytes
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:27Z
---
Mirela won't let a command reach Pipit-2 unless it carries a signature only Apogee can make: an HMAC, computed from the command and a secret key that Apogee and the satellite share.

`hmac.New(sha256.New, key)`, from `crypto/hmac` and `crypto/sha256`, gives a hash you write bytes into; `Sum(nil)` returns the 32-byte result:

```go
mac := hmac.New(sha256.New, []byte("demo key"))
mac.Write([]byte("furrow-a|7|PING"))
fmt.Printf("%x\n", mac.Sum(nil)[:4])
// 32cb6b96: the first 4 of 32 bytes
```

To check a signature, compute your own and compare the two with `hmac.Equal`. Unlike `bytes.Equal`, it takes the same time wherever the bytes differ, so timing gives an attacker no clue.

## Your task

1. Make `sign` return the HMAC-SHA256 of `canonical` under `key`, in lowercase hex.
2. Make `verify` decode `sigHex`, return `false` if it isn't hex, and compare it with the HMAC it computes, using `hmac.Equal`.
