---
slug: ready-for-bigger-computers
title: Ready for Bigger Computers
type: code_tests
position: 6
difficulty: 8
topics: crypto, security
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
Wednesday evening, Tomás leans over your desk with a printout: "Record the backhaul tonight and keep it. If anyone ever builds a big enough quantum computer, a machine built on quantum physics, it could work out X25519's secret from the two shares. Recordings keep!" He wants a key exchange made for that day.

ML-KEM-768 is designed to resist such machines, which is what post-quantum means. It works one way round: the receiver publishes an encapsulation key of 1,184 bytes, and the sender makes a new 32-byte key and a 1,088-byte ciphertext that only the receiver's decapsulation key opens:

```go
dk, _ := mlkem.GenerateKey768() // receiver
ek := dk.EncapsulationKey().Bytes()
fmt.Println(len(ek)) // 1184
```

## Your task

1. Make `wrapFor` read `ek` with `mlkem.NewEncapsulationKey768` and return the key and ciphertext from its `Encapsulate`.
2. Make `unwrap` return what `dk.Decapsulate` makes of `ciphertext`.

> [!NOTE]
> A ciphertext changed on the way still opens, to a different key and no error. The session then fails at its first sealed message, which the receiver can't open.
