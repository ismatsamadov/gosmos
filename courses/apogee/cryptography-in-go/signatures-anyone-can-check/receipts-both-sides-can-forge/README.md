---
slug: receipts-both-sides-can-forge
title: Receipts Both Sides Can Forge
type: code_tests
position: 0
difficulty: 8
topics: crypto, security
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Monday 28 August, 09:02. Hollis: "FURROW-A 0827 1402 NEVER ARRIVED." The first draft of receipts, signed notes that a file was delivered, signs with a secret Furrowlight holds too, so the test, playing Furrowlight, signs one saying the file landed a day late.

Ed25519 signs with a private key and checks with the matching public key, 32 bytes anyone may hold. Whoever has only the public key can check a signature but never make one:

```go
pub, priv, _ := ed25519.GenerateKey(nil)
msg := []byte("PING")
sig := ed25519.Sign(priv, msg)
fmt.Println(len(pub), len(sig)) // 32 64
fmt.Println(ed25519.Verify(pub, msg, sig)) // true
```

## Your task

1. Make `Notary` hold an `ed25519.PrivateKey` from `ed25519.NewKeyFromSeed`. `newNotary` returns an error for a seed that isn't 32 bytes, where that function would panic.
2. Make `Sign` sign with it, and `Published` return its public key.
3. Make `checkReceipt` verify with `ed25519.Verify`, returning false for a key that isn't 32 bytes, where `Verify` would panic.
