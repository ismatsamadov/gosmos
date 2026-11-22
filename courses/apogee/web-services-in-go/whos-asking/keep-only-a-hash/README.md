---
slug: keep-only-a-hash
title: Keep Only a Hash
type: code_tests
position: 2
difficulty: 8
topics: security, strings
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Mirela, on the customer table: "If a backup of this leaks, nobody should be able to book with what is in it." So from today the service keeps only each key's hash.

SHA-256 turns any text into 32 bytes, its hash. The same text always gives the same 32 bytes, and the only way back to the text is guessing it, hopeless against a key's 32 random hex digits. `sha256.Sum256` returns them as a `[32]byte` array, and `hex.EncodeToString` writes bytes as hex digits, two for each byte, in lowercase, as `KeyHash` holds them:

```go
sum := sha256.Sum256([]byte("Port Ellery"))
fmt.Println(len(sum)) // 32
fmt.Println(hex.EncodeToString(sum[:4]))
// 7009b310
```

## Your task

The customer table now holds `KeyHash`, the hash of each customer's whole key, and `checkKey` still compares the key itself with it, so it refuses everyone. Hash the whole key with `sha256.Sum256`, hex-encode the sum, and compare that with the customer's `KeyHash` in constant time.

> [!TIP]
> Slice the array, `sum[:]`, to pass it where a `[]byte` is needed.
