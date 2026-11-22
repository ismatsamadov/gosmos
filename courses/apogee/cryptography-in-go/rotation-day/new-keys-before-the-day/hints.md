`satKey(next.Key, "sign", satID)` and the same with `"seal"` give the two 32-byte keys, and `hex.EncodeToString` writes each as the 64 hex digits the plaintext carries.

The command is a `Command` with `Op` `"LOAD_KEY"` and three `Args`: `next.ID`, the sign key's hex, then the seal key's. `sealCommand` joins them with `|` itself and adds a fresh nonce.

If Cirrine-1 answers `unknown key`, the line is sealed under `k-2028-10`, which it doesn't hold yet. Ask `signingKey` which version seals the second before `next.Activate`, and seal under that version's `ID`.
