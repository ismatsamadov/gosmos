`ed25519.NewKeyFromSeed(seed)` returns the private key: keep it in a field of type `ed25519.PrivateKey` in place of `secret`. Compare `len(seed)` with `ed25519.SeedSize` first, since the function panics on any other length.

`ed25519.Sign` takes the private key and the message. A private key knows its public half: its `Public` method returns it as an interface value, and a type assertion to `ed25519.PublicKey` gives the 32 bytes.

In `checkReceipt`, return false when `len(published)` isn't `ed25519.PublicKeySize`, else what `ed25519.Verify` says. Once nothing uses HMAC, remove `crypto/hmac` and `crypto/sha256` from the imports.
