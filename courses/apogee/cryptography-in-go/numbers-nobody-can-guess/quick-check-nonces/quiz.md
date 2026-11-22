## What is true of a random 12-byte nonce under one key?

- [ ] It must stay as secret as the key
  > Anyone listening to the uplink may read it; what keeps a command safe is the key.
- [x] It may be public, and one key can use 2^32^ random ones safely
  > A nonce must only never repeat under its key, and with 12 random bytes the odds of any repeat in 2^32^ commands stay below one in eight billion.
- [ ] It is safe only when counted up from 0
  > Counting is what failed: a count kept in memory starts at 0 again after a restart.

## After uplinkd restarts, what does a nonce counter kept in memory hand out?

- [x] 0, 1, 2 and so on again: nonces the key has already used
  > The count lived in memory, which the restart cleared, so the new run repeats the old run's first nonces under the same key.
- [ ] The number after the last nonce it handed out
  > Nothing remembers that number: it lived in the memory the restart cleared.
- [ ] Nothing, until someone sets the count by hand
  > The counter starts again at 0 by itself, which is exactly the trouble.
