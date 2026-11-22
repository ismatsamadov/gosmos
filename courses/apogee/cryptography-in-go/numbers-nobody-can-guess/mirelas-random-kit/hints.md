Three of the helpers follow one pattern, the one from Monday's nonce: make a slice of the size the table gives, fill it with `rand.Read`, and return it. For `newKey`, the filled slice becomes the key's `Secret`.

`newSessionToken` needs no slice: `rand.Text` already returns a whole token. Once every helper calls `rand.Read` or `rand.Text`, remove the `_ = rand.Read` line.

If `TestOnlyCryptoRand` fails, read its message: it names the file and the import it found. If a size test fails, compare your slice's length with the table: 32 for a key, 12 for a nonce, 16 for a salt.
