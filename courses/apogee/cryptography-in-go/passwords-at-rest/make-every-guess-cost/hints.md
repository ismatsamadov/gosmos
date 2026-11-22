Start with the constant: `iterations` is the count every new record will use, so it holds 600000 rather than 1.

Replace the body of `derive` with one call to `pbkdf2.Key`, which already returns a key and an error, so you can return both as they come. Add `"crypto/pbkdf2"` to the imports.

If `TestDerive` still fails, check the arguments' order: the hash constructor `sha256.New` without parentheses, then the password, the salt, the count and a length of 32.
