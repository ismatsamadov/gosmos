## What does a salt stop on its own?

- [ ] Every guess, because the attacker can't read the salt
  > The salt sits in the record beside the key, so anyone holding the file reads it; a single guess costs no more than before.
- [x] Tables made in advance, and equal passwords storing the same
  > Each record gets its own salt, so a table of digests made beforehand matches nothing and every record must be attacked alone; each guess is still as cheap as one hash.
- [ ] Nothing, because the attacker has the salt too
  > Even in the open, the salt makes every record different, which is what defeats a table worked out before it existed.

## What do 600,000 iterations cost someone holding the user file?

- [ ] One hash per guess, since only the console runs the iterations
  > To test a guess against a record, the attacker must derive the key exactly as the console does, all 600,000 rounds of it.
- [ ] Nothing, once they have the salt and the count from the record
  > Knowing the count tells them how much work to do; it doesn't let them skip any of it.
- [x] 600,000 rounds of work for every guess at every record
  > There is no shortcut through PBKDF2, so a billion guesses a second become fewer than a thousand, for each record separately.

## Must the salt be kept secret?

- [x] No: it is stored in the record, and checking a password needs it
  > The salt only has to be new for every record; the console reads it back from the record whenever someone logs in.
- [ ] Yes: whoever knows the salt can work out the password
  > The salt and the key together still leave the password to be guessed, one 600,000-round derivation at a time.
- [ ] Yes, so it belongs in a separate file from the keys
  > A second file would add nothing: the salt's job is to be different for every record, not to be hidden.
