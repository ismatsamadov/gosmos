`cipher.NewGCM` takes the block that `aes.NewCipher` already made and returns a `cipher.AEAD` and an error. Make it in both functions, right after the block, and return early if it fails, as the block's error does.

The AEAD's `Seal` takes four slices: where to append the result (`nil` for a new slice), the nonce, the plaintext and the associated data; `Open` takes the same, with the sealed bytes in the plaintext's place, and returns the plaintext and an error. `counter` and the calls to `NewCTR` can go.

If `TestBitFlip` still gets `SET_MODE|LEOP` back, a function still uses counter mode: both must seal and open with the AEAD, since bytes sealed by one mode never open in the other.
