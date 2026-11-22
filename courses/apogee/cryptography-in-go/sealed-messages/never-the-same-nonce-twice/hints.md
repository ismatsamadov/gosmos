`random.go` already has `newNonce`, which returns 12 bytes from `crypto/rand`. Call it where the draft decodes its fixed hex.

`nonce` is a named result: give it the value `newNonce` returns, and both the return and `Seal` use it. Once nothing decodes hex, remove the `encoding/hex` import.

If `TestMirelaReads` still reads `DUMP_FILE|occ`, two seals still share a nonce: take the new one inside `sealCmd`, once per call, never from a constant or a variable kept outside it.
