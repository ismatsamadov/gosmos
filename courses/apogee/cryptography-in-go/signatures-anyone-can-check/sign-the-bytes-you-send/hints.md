`e.Receipt` holds the receipt's bytes exactly as they arrived. Hand those, not a new encoding, to `checkSig` with `pub` and `e.Sig`, and return its error when it has one.

Only after the signature checks, decode `e.Receipt` into a `Receipt` with `decodeStrict`, passing its address, and wrap its error with `ErrBadEnvelope` the way the envelope's own decoding error is wrapped.

If `TestRefund` reports `ErrBadEnvelope` where it wants `ErrBadSignature`, something still decodes the receipt before the signature is checked. Nothing needs `json.Marshal` anymore: the receipt's bytes are already the ones Apogee signed.
