`sha256.Sum256` takes the whole file and returns its 32-byte digest as an array; the CRC-32 line has nothing left to do.

`fmt.Sprintf` with `%x` writes every byte of the array as two hex digits, 64 in all. `hex.EncodeToString` does the same with a slice of the array.

If the report's fingerprint starts `f6ced6c0`, the CRC-32 is still returned. If it has 8 digits only, check that you print the whole digest, not part of it.
