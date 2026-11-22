`decodeEnvelope(env)` decodes into `Envelope` itself and refuses any other field, so an envelope that brings `pub` goes no further. The struct with `Pub` can go, and `encoding/hex` and `encoding/json` with it.

Looking up `e.KeyID` in `keys` with `ok` finds the key Apogee publishes under that ID. When `ok` is false, return `ErrUnknownSigner` wrapped with the ID; otherwise `openReceipt`, given that key and the envelope, checks the signature and decodes the receipt.

Compare the receipt's `KeyID` with `e.KeyID` before returning it, and wrap `ErrBadEnvelope` when they differ. If `TestKeyItBrings` still fails, the key handed to `openReceipt` comes from somewhere other than `keys`.
