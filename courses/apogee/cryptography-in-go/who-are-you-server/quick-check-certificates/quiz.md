## What does a certificate bind together?

- [ ] A name and a password, which the server checks
  > A certificate holds no secret at all: anybody may read it, so a password in it would leak.
- [x] A public key and names, signed by its issuer
  > The issuer's signature vouches that this public key belongs to these names, for these dates and uses.
- [ ] A private key and a name, so its holder can prove itself
  > The private key never goes into a certificate: its holder keeps it, and the certificate carries the public half.

## A certificate copies every one of Ystmark's names. Why does `Verify` still refuse it?

- [ ] Its names are spelled a little differently
  > They can be spelled exactly the same: names cost nothing to copy.
- [ ] Its dates fall before Apogee's root was made
  > Its maker can choose any dates; the dates aren't what give it away.
- [x] No chain of signatures leads from it to Apogee's root key
  > Anybody can write any name, but only Apogee's root key, or the stations CA's key that the root certified, makes a signature that chains to Apogee's root.

## Why does the root's private key stay in the safe?

- [x] Whoever holds it can certify any name at all
  > Every certificate it signs chains to Apogee's root, so its holder could make a valid certificate for any station or server.
- [ ] Copying the root's certificate would let anybody sign with it
  > The root's certificate is public and copying it changes nothing; only the private key signs.
- [ ] A key wears out when it is used on a network
  > A key doesn't wear out; it stays offline so that nobody can steal it.
