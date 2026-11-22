Both sides make a key with `ecdh.X25519().GenerateKey` and send, with `writeLine`, the line `hello` makes from `me` and that key. `dial` writes first, then calls `readHello`; `answer` calls `readHello` first, then writes.

After `readHello` returns the other end's share, `mine.ECDH(theirs)` gives the secret, and `sessionKey` takes it with the dialer's share first: in `dial` that is your own `mine.PublicKey().Bytes()`, in `answer` it is `theirs.Bytes()`. In `dial`, compare the code `readHello` returns with `peer` before deriving anything.

If the test reports two different keys, the shares went into `sessionKey` in different orders at the two ends. If it reports an i/o timeout after about a second, both ends wrote first, or both read first.
