Call `newSalt` once, at the top of `hashPassword`, so every call gets its own 16 bytes. A hash from `sha256.New()` takes the salt and then the password, one `Write` after the other, and `Sum(nil)` returns the digest.

The record has three fields joined by `$`: the word `salted`, the salt in hex and the digest in hex. `hex.EncodeToString` writes both in lowercase, two digits per byte.

If the digest doesn't recompute, check the order: the salt's bytes come first, then the password's. If two records for one password are still equal, the salt isn't new on every call.
