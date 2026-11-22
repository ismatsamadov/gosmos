---
slug: thirty-two-random-bytes
title: Thirty-Two Random Bytes
type: code_tests
position: 1
difficulty: 6
topics: crypto, errors
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
On Friday, Mirela makes the secret that every satellite's uplink keys will be made from, the season's master key, `k-2028-07`. Before then she wants the helper that makes it, and the parser that reads it back from the key file.

A key of 256 bits is 32 random bytes, one of 2^256^ possible keys. In the key file it is 64 hex digits, two per byte. `hex.DecodeString` reads hex, but it takes any even number of digits, even none:

```go
b, err := hex.DecodeString("c0ffee")
fmt.Println(len(b), err) // 3 <nil>
b, err = hex.DecodeString("")
fmt.Println(len(b), err) // 0 <nil>
```

## Your task

1. Make `newKey` fill the key's `Secret` with 32 bytes from `rand.Read`.
2. Make `parseSecret` return an error unless `text` is exactly 64 hex digits, so a line cut short never loads as a 31-byte key.

> [!NOTE]
> The tests' master key is 32 readable characters, `uplink test master key, 32 bytes`. Test keys are readable on purpose; flight keys never appear in code.
