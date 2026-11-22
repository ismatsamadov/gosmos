---
slug: shares-on-the-wire
title: Shares on the Wire
type: code_tests
position: 1
difficulty: 8
topics: crypto, errors
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
Ystmark's share reaches Port Ellery as one line of 64 hex digits, over a backhaul that garbles a line now and then. Sunniva: "Bad line, refuse it. No guessing."

A share has exactly one spelling, 64 lowercase hex digits, so the text that carries it can be signed and compared later. A well-formed share can still be degenerate: 32 zero bytes bring every private key to an all-zero secret, and since `NewPublicKey` checks only the length, `ECDH` is what refuses it, with an error. `hex.DecodeString` is looser than one spelling: capitals decode to the same bytes.

```go
a, _ := hex.DecodeString("c0ffee")
b, _ := hex.DecodeString("C0FFEE")
fmt.Println(bytes.Equal(a, b)) // true
```

## Your task

1. Make `parseShare` return an error unless `text` is exactly 64 lowercase hex digits. Sixty-four zeros are well formed, so they parse.
2. Make `agreeHex` return the error from `ECDH` instead of dropping it, so a degenerate share stops the agreement.

> [!WARNING]
> Dropping that error leaves the session an empty secret, and a key made from it is one anybody can make.
