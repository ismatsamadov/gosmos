## Someone reads Pipit-2's key out of the spare computer on the bench. What else can they work out?

- [ ] The master key, by running the derivation backward
  > HKDF can't be run backward: from a derived key there is no way to the master.
- [x] Nothing more: not the master, and no other satellite's key
  > Every derived key is unrelated to the others, so the stolen key commands Pipit-2 and nothing else.
- [ ] Every other satellite's key, since all of them share one master
  > Without the master itself, nobody can derive another satellite's key, however many derived keys they hold.

## Why does a password never go through HKDF?

- [ ] HKDF refuses secrets shorter than 32 bytes
  > HKDF takes a secret of any length; the trouble is how cheaply a password can be guessed.
- [ ] HKDF can only make keys longer than any password
  > The key is as long as you ask for; length has nothing to do with it.
- [x] HKDF is fast, so every guess at the password would cost almost nothing
  > HKDF is for secrets that are already random; a password needs PBKDF2's slow work, 600,000 rounds a guess.

## Mirela changes the salt's `v1` to `v2`. What happens to the derived keys?

- [x] Every one changes, for every satellite and purpose
  > The salt goes into every derivation, so all of them change, and no satellite would accept the new keys until it holds them.
- [ ] Only the sign keys change, since the salt was made for them
  > The salt is the same for every derivation, so the keys for every purpose change with it.
- [ ] None of them, because the salt is public anyway
  > Public isn't unused: the salt is mixed into every byte HKDF returns.
