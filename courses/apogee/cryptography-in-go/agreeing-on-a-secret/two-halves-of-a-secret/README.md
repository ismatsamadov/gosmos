---
slug: two-halves-of-a-secret
title: Two Halves of a Secret
type: code_tests
position: 0
difficulty: 8
topics: crypto, security
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
Monday 4 September. Since Sunday's storm cut Ystmark's fiber cable to the mainland, the station reaches headquarters at Port Ellery only over a slow satellite link another company runs, the backhaul. Mirela: "This week the two ends need a key nobody on that link can learn, and neither end can send it."

Two ends can still agree on a key in the open. Each makes a private key and sends only its public half, called its share. Your private key with their share gives the same 32 bytes as their private key with yours, and with today's computers nobody holding only the two shares can work them out: a key agreement. Go's X25519 makes the keys:

```go
curve := ecdh.X25519()
priv, _ := curve.GenerateKey(rand.Reader)
share := priv.PublicKey().Bytes()
fmt.Println(len(share)) // 32
```

## Your task

`agree` gets your private key and the other end's share as bytes. Make it read the share with `ecdh.X25519().NewPublicKey`, then return `mine.ECDH` of it, passing on any error either one returns.
