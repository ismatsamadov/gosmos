`rsa.EncryptOAEP` takes the hash first, as `sha256.New()`, then `rand.Reader`, the public key, the message and the label as a `[]byte`. What it returns is exactly what `encryptTo` returns.

`rsa.DecryptOAEP` takes the same kind of hash, a source of randomness it no longer uses (`nil` is fine), the private key, the ciphertext and the same label.

If Brinewatch can't open what `encryptTo` locked, though `decryptWith` can, the two sides differ in hash or label: both must be SHA-256 and `[]byte(label)`. If `TestLabelBinds` fails, one side drops the label.
