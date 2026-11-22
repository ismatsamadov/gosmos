`x509.MarshalPKIXPublicKey` returns the key's standard bytes and an error. Put the bytes in a `pem.Block` whose `Type` is `PUBLIC KEY`, and give its address to `pem.EncodeToMemory`.

`pem.Decode` returns a nil block when it finds none. Check the block's `Type`, and that `bytes.TrimSpace(rest)` is empty, before parsing `block.Bytes` with `x509.ParsePKIXPublicKey`.

`ParsePKIXPublicKey` returns an `any`, and a type assertion to `ed25519.PublicKey` with `ok` keeps the key only if it is an Ed25519 one. When `ok` is false, return `errors.New("not an Ed25519 key")`, with exactly that text.
