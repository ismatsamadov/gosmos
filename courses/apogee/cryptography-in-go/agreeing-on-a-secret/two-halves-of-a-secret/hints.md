`ecdh.X25519()` returns the curve, and its `NewPublicKey` turns 32 bytes into an `*ecdh.PublicKey`, or returns an error for any other length.

If `NewPublicKey` returns an error, return it with a nil secret. Otherwise `mine.ECDH` takes the public key and returns the secret and an error, which `agree` can return as they come.

If `TestBothWays` still reports 0 bytes, the placeholder line and `return nil, nil` are still there: replace both with the parse and the call to `ECDH`.
