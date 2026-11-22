Replace the type assertion with `switch k := key.(type)`. Inside each `case`, `k` has that case's type, so in the `*rsa.PublicKey` case `k.N.BitLen()` counts the modulus's bits.

Ed25519 public keys arrive as `ed25519.PublicKey`, not a pointer, and X25519 ones as `*ecdh.PublicKey`: import `crypto/ed25519` and `crypto/ecdh` for those two cases. In `default`, `fmt.Errorf` with `%T` names whatever type came.

If Brinewatch's 3,072-bit keys are refused, or the key of 2,047 bits gets through, the check counts bytes with `Size()`, where `N.BitLen()` counts bits. If the impostor's key of exactly 2,048 bits is refused, the comparison needs `<`, not `<=`.
