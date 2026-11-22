`seed` is an array of 32 bytes, and `ed25519.NewKeyFromSeed` takes a slice, so slice the array whole. The private key it returns gives its public half through its `Public` method and a type assertion to `ed25519.PublicKey`, as your notary's `Published` does.

`ed25519.Sign` takes the private key and the message and returns the 64-byte signature; `ed25519.Verify` takes the public key, the message and the signature, and reports whether they fit. For the last line, verify a `bytes.Clone` of the receipt with one byte XORed, against the same signature.

The first line prints all 32 bytes until you slice them: `pub[:8]` gives the first 8, which `%x` writes as 16 hex digits. If they aren't `ba6c6226d2c3c056`, check that the seed is the SHA-256 of exactly `receipts r-2028-08`.
