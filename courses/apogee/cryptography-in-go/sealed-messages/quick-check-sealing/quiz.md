## Someone changes one byte of `sealed` on its way to the satellite. What does `Open` return?

- [ ] The command, with one byte changed
  > That is what counter mode alone gives, as in Monday's draft; GCM's tag covers every sealed byte.
- [x] An error, and none of the command
  > The tag no longer matches, so `Open` refuses before returning a single byte: nothing tampered is ever read.
- [ ] The command as it was, repaired with the tag
  > A tag only shows that something changed; it holds nothing to repair the bytes with.

## Where may a command's nonce travel?

- [x] In the open beside the command, as long as it never repeats under that key
  > A nonce isn't secret: `Open` needs it, so it goes up in the clear, and only a repeat under one key gives commands away.
- [ ] Only inside the sealed bytes, where nobody can read it
  > `Open` needs the nonce before it can open anything, so it can't travel locked inside the bytes it opens.
- [ ] Nowhere: both ends work it out from the time of the pass
  > Every command sealed for that pass would get the same nonce, a repeat under one key, so each nonce comes from `crypto/rand` and goes up beside its command.

## What does `Open` return with the right key and the wrong nonce?

- [ ] The command, because the key is right
  > The tag depends on the nonce as well as the key, so another nonce gives another tag.
- [ ] Scrambled bytes, and no error
  > That is what counter mode alone would give; GCM checks the tag before returning anything.
- [x] An error, as with a wrong key
  > The tag is worked out from the key, the nonce and every sealed byte, so the wrong nonce fails the check.
