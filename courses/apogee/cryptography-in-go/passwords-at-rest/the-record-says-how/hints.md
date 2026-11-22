`formatRecord` is one `fmt.Sprintf`: `%d` writes the count, and `%x` writes a byte slice as lowercase hex, two digits per byte.

In `parseRecord`, split at every `$` and require exactly four fields, the first being `pbkdf2-sha256`, before you read any of them. Then `strconv.Atoi` reads the count, and `hex.DecodeString` the salt and the key.

If a refusal still gets through, check every field after decoding it: the count between 1000 and 10000000, the salt exactly 16 bytes and the key exactly 32. A decode that succeeds can still hold the wrong number of bytes.
