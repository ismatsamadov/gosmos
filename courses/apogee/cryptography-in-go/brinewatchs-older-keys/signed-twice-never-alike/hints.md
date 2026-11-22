`sha256.Sum256(req)` returns an array, and `[:]` makes it the slice both functions sign and check. `crypto.SHA256`, from the package `crypto`, tells RSA which hash made it.

`rsa.SignPSS` takes `rand.Reader`, the private key, the hash, the digest and the options. `rsa.VerifyPSS` takes the public key, the hash, the digest, the signature and the options, and returns nil only for a good signature.

If signing fails with `input must be hashed with given hash`, the request itself went where its 32-byte digest belongs. If a changed request is accepted, `checkRequest` isn't returning what `VerifyPSS` says.
