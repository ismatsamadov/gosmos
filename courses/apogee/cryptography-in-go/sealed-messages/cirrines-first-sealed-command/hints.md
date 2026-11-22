`canonical` refuses a command without one single meaning, and `ring.SatKey` refuses a key ID the ring doesn't hold or a satellite outside the fleet. Return `""` and the error as soon as either fails, before anything is sealed.

`newAEAD(key)` gives the AEAD. Build the header first, a `Sealed` with version 2, the command's satellite and number, `keyID` and `Nonce: newNonce()`, since `s.AD()` reads it. The plaintext is `strings.Join` of the operation followed by the arguments, with `|` between them.

Set `s.Ciphertext` to what `aead.Seal` returns for the nonce, the plaintext and `s.AD()`, then return `formatSealed(s)`. If the stand-in answers `ErrForged`, the associated data differs from what it opens with: seal with `s.AD()` of the very header you format. If it opens another command than yours, the plaintext must be `OP|ARG`, without the satellite or the number.
