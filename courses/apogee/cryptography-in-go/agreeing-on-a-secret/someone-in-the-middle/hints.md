In `hello`, the bytes to sign are the readme's `signed:` line: the words `backhaul share`, the code and the share's hex, joined by `|`. Sign them with `ed25519.Sign` and `id.Key`, and write the signature's hex after the share, one space apart.

In `checkHello`, check in this order: four fields starting with `HELLO`; the code's key in `known`; the share through `parseShare`; then decode the signature and `ed25519.Verify` it with that key over the same bytes `hello` signs. Each refusal is `fmt.Errorf("%w: ...", ErrBadHello, ...)`.

If the middle's HELLO still gets through, `checkHello` verifies with the wrong key or the wrong bytes: the key is `known[code]`, and the bytes are `backhaul share|`, the code, `|` and the share's hex exactly as it arrived.
