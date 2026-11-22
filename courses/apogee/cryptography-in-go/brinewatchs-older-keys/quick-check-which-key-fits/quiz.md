## Port Ellery and Ystmark need a fresh secret for every backhaul session. Which fits?

- [ ] RSA, locking each session's secret to Ystmark's lasting key
  > It would work today, but whoever steals that lasting key later opens every session recorded before.
- [x] X25519, with a new key pair each session
  > Each end makes a key for the session and forgets it afterwards, so a key stolen later opens nothing recorded.
- [ ] Ed25519, signing each session's secret
  > A signature proves who sent something; it hides nothing and agrees on no secret.

## Why do Apogee's receipts stay Ed25519 rather than switching to RSA?

- [x] Ed25519 is as strong with 64-byte signatures and 32-byte public keys
  > A 3,072-bit RSA key gives about the same strength with 384-byte signatures and keys twelve times the size.
- [ ] RSA keys can't sign anything
  > They can: a PSS signature is an RSA signature, as Brinewatch's security office makes them.
- [ ] Ed25519 keys never need replacing
  > Every key is replaced in time: `r-2028-08` stops signing on Rotation Day.

## Why does Brinewatch's handover use RSA at all?

- [ ] RSA is newer than Ed25519 and X25519
  > RSA is the oldest of the three by decades, which is why so many older systems speak it.
- [ ] Only RSA can keep a secret readable by one company
  > X25519 with a cipher could do it too, if Brinewatch's systems spoke it.
- [x] Brinewatch's systems speak only RSA, and RSA holds up when used with care
  > With OAEP, PSS and 3,072 bits, RSA holds up against today's computers, only bigger and slower than the newer kinds.
