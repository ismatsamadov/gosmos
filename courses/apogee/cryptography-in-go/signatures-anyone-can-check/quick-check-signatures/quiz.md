## What can a customer do with Apogee's public receipt key?

- [x] Check Apogee's signatures, but never make one
  > Checking needs only the public key; signing needs the private key, which never leaves bookingd.
- [ ] Check Apogee's signatures, and make ones of its own that Apogee's check accepts
  > A signature that verifies under Apogee's public key can only come from Apogee's private key.
- [ ] Nothing, until Apogee sends the private key too
  > The public key alone checks every signature; the private key is never sent anywhere.

## bookingd signs the same receipt twice with r-2028-08. What are the two signatures?

- [ ] Different, since every signature takes fresh random bytes
  > Some signing schemes work that way, but Ed25519 uses no randomness when it signs.
- [x] The same 64 bytes both times
  > Ed25519 signs the same message with the same key the same way every time, which is why a test may expect a fixed signature.
- [ ] Different, since a signature includes the time it was made
  > Only the message is signed; a receipt carries its own `signed_at`, and the same bytes sign the same way.

## Why can't a signature made with a secret both sides hold settle Hollis's complaint?

- [ ] HMAC-SHA256 is too weak to trust
  > HMAC-SHA256 holds up well; the trouble is who holds its key, not its strength.
- [ ] It covers the message about the file, not the file itself
  > Even a signature over every byte of the file would prove nothing if either side could make it.
- [x] Either side could have made it, since both hold the key
  > With a shared secret, Furrowlight could make Apogee's signature and Apogee Furrowlight's, so it proves nothing to anyone else.
