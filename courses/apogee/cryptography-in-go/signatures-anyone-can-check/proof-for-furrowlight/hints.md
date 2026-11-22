Each check is a call and an `if`: `decodeEnvelope`, then `findKey` (wrap `ErrUnknownSigner` with the ID when it reports false), then the key's `Revoked`, then `openReceipt` with the key's `Public`, then the receipt's `KeyID` against the envelope's. Return as soon as one fails.

The window refuses a `SignedAt` before `k.Activate`, and one at or after `k.Retire` unless `Retire` is zero; `Before` and `IsZero` answer both questions. Only then call `checkFile`.

In `checkFile`, `parseDigest(r.Digest)` gives the hash's name and the digest's bytes, and `newHash` the hash itself. `io.Copy(h, file)` returns how many bytes it hashed: compare that with `r.Bytes`, and `h.Sum(nil)` with the digest's bytes using `bytes.Equal`.
