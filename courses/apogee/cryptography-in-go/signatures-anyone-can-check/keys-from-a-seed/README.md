---
slug: keys-from-a-seed
title: Keys From a Seed
type: code
position: 1
difficulty: 7
topics: crypto, printing
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Mirela makes the receipt key for the tests from a seed, the SHA-256 of `receipts r-2028-08`, so that every run of every test signs the same way.

An Ed25519 private key is made from 32 bytes called its seed. `ed25519.NewKeyFromSeed` makes the same key from the same seed, and Ed25519 signs a message the same way every time, so a test can expect a fixed signature. A real seed is 32 random bytes; a test's is readable on purpose. `sha256.Sum256` returns an array, which `[:]` turns into a slice:

```go
sum := sha256.Sum256([]byte("PING"))
fmt.Println(len(sum), len(sum[:4])) // 32 4
```

## Your task

`main` prints placeholders. Make the key from `seed`, sign `receipt` with it, and print the public key's size and first 8 bytes, the signature's size, whether it verifies, and whether it verifies with one byte of the receipt changed:

```text
public key: 32 bytes, ba6c6226d2c3c056
signature: 64 bytes
verifies: true
one byte changed: false
```
