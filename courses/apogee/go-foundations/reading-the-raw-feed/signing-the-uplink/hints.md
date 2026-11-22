Both functions need the same HMAC of `canonical` under `key`. In the readme's example, where do the key and the command go in?

For `sign`, `hex.EncodeToString` from `encoding/hex` turns the 32 bytes of `Sum(nil)` into the 64 hex digits the wire format carries.

For `verify`, decode `sigHex` with `hex.DecodeString` and return `false` on its error; then pass the decoded bytes and your own HMAC to `hmac.Equal`.
