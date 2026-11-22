`sha256.Size` is the constant 32, the length of every SHA-256 digest: print it where `crc32.Size` is now.

`sha256.Sum256` returns an array. Slice its first 8 bytes with `[:8]` and print them with `%x`, which writes two hex digits for every byte.

If a line shows 64 hex digits, you printed the whole digest rather than its first 8 bytes. If the build stops at `"hash/crc32" imported and not used`, remove that import.
