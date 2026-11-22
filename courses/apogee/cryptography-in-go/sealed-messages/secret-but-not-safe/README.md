---
slug: secret-but-not-safe
title: Secret but Not Safe
type: code_tests
position: 0
difficulty: 8
topics: crypto, security
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Monday 14 August. Mirela's first draft locks Cirrine-1's commands with AES, the standard cipher, in counter mode. Without the key, she turns a locked `SET_MODE|SAFE` into `SET_MODE|LEOP`, the mode of a satellite's first days after launch.

Counter mode XORs the command with a stream of bytes made from the key and the nonce, so flipping a bit of the locked bytes flips the same bit of the command:

```go
ks := byte(0x5a)         // a byte of the stream
c := 'S' ^ ks            // the locked byte
c ^= 'S' ^ 'L'           // flipped, no key
fmt.Printf("%c\n", c^ks) // L
```

GCM, Galois/Counter Mode, adds a 16-byte tag worked out from the key, the nonce and every locked byte, and its `Open` refuses bytes the tag doesn't match. A command locked so that nobody can read it or change it unnoticed is a sealed command.

## Your task

Make `seal` and `open` use GCM: wrap the AES block with `cipher.NewGCM`, then call its `Seal` and `Open`, with `nil` as the last argument for now. The test flips Mirela's bits.
