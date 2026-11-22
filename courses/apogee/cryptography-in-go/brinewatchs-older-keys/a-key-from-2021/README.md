---
slug: a-key-from-2021
title: A Key From 2021
type: code_tests
position: 3
difficulty: 9
topics: crypto, interfaces, switch
created: 2026-11-13T03:46:39Z
modified: 2026-11-13T12:17:21Z
---
Friday, 16:10. With the pad retired, Brinewatch's security office emails the public key its systems have used to receive secrets since 2021, as a PEM block. Mirela counts its bits before she reads anything else: 1,024.

Brinewatch's keys are of an older kind, called RSA: the public key checks signatures and can also lock a message that only the private key opens. Its strength grows with the size of its modulus, `N`, a number thousands of bits long. Keys of 1,024 bits are no longer counted safe, and Apogee takes 2,048 and up. `x509.ParsePKIXPublicKey` returns an `any`, and a type switch asks which kind of key is inside:

```go
switch v := x.(type) {
case string:
	fmt.Println("text of", len(v), "bytes")
case int:
	fmt.Println("the number", v)
}
```

## Your task

Make `readRSA` switch on the key's type. Return an `*rsa.PublicKey` whose `N.BitLen()` is 2048 or more. Refuse a smaller one, an `ed25519.PublicKey` and an X25519 `*ecdh.PublicKey`, and name any other type with `%T`:

```text
want at least 2048 bits, got 1024
want an RSA key, got Ed25519
want an RSA key, got X25519
want an RSA key, got *ecdsa.PublicKey
```
