## The whiteboard's code runs again tomorrow. What key does it make?

- [ ] A new key, since a day has passed
  > The generator's numbers depend only on its seed, and the seed is 7 and 7 on every run.
- [x] The same 32 bytes as today
  > A generator started from a known seed gives the same numbers on every run, so anyone with the seed has the key: no secret comes from `math/rand/v2`.
- [ ] A key that only this Mac can make again
  > Any computer with Go and the same seed makes exactly the same bytes.

## What does `rand.Read` from `crypto/rand` promise about the slice you give it?

- [ ] It fills part of it, and the count it returns says how much
  > An `io.Reader` may stop short like that, but `rand.Read` always fills the whole slice.
- [ ] It fills it unless it returns an error, which you must check first
  > It never returns an error: if the Mac can't supply randomness, the program stops instead.
- [x] It fills every byte, or the program stops
  > It never returns an error and always fills the slice; without randomness to give, the program stops, so code calls it without checking.

## How many different values can a 32-byte key take?

- [x] 2^256^, because 32 bytes are 256 bits
  > Every bit doubles the count, so 256 bits give 2^256^ possible keys.
- [ ] 32 × 256, or 8,192
  > That adds the choices byte by byte; each byte multiplies the count by 256 instead.
- [ ] 2^32^, about four billion
  > That counts each byte as a single bit; a byte is 8 bits, so 32 bytes are 256.
