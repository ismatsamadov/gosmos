## Why does the check compare with `subtle.ConstantTimeCompare` and not `==`?

- [ ] `==` can't compare two strings of 64 hex digits
  > `==` compares strings of any length; the trouble is how long it takes, not whether it works.
- [x] `==` can stop at the first difference, so its timing tells a guesser how much matched
  > The more that matches before the first difference, the later the refusal, a clue a patient guesser can use; `ConstantTimeCompare` takes the same time wherever the bytes differ.
- [ ] `ConstantTimeCompare` is faster for long keys
  > It isn't faster: it keeps reading after a difference, and that is the point.

## A backup of the customer table leaks, every `KeyHash` in it. What can the thief book with them?

- [ ] Any customer's passes, by sending a `KeyHash` where the key goes
  > The service hashes whatever it is sent, and the hash of a `KeyHash` is a different hash.
- [ ] Any customer's passes, once the keys are worked back out of the hashes
  > SHA-256 has no way back: the only road from a hash to its key is guessing keys, and 32 random hex digits are far too many to guess.
- [x] Nothing: a hash isn't a key, and finding the key from its hash is hopeless
  > The service accepts only a key itself, and a stolen hash gives the thief no shorter road to one than guessing.

## Brinewatch's valid token asks to book a pass of `cirrine-1`. Which status fits?

- [x] `403 Forbidden`
  > 403 says "we know who you are, and you may not": the token names Brinewatch, and the satellite is Cirrine's.
- [ ] `401 Unauthorized`
  > 401 says "we don't know who you are", and this token says exactly who Brinewatch is.
- [ ] `400 Bad Request`
  > The request is well formed; what's wrong is who is asking.
