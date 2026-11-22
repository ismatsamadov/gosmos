Make the hash with `sha256.New()` and pass it to `io.Copy` as the writer, with `r` as the reader: the copy reads until the stream ends and writes each piece into the hash.

`io.Copy` returns a count and an error. When the error isn't nil, return `""` with it; otherwise `h.Sum(nil)` is the digest, and `hex.EncodeToString` turns it into 64 hex digits.

If the 1 MiB stream still gets the wrong digest, the old `Read` is still in charge: remove the buffer and the single `Read`, and let `io.Copy` do all the reading.
