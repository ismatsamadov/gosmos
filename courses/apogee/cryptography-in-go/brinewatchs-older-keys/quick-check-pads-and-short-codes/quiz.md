## Someone kept both emails. What do they get by XORing the two together?

- [ ] Nothing, without a copy of the pad
  > The same pad bytes are in both, so XORing the emails cancels them: no copy of the pad is needed.
- [x] The two secrets XORed together, so knowing one reveals the other
  > With the pad cancelled, what's left is the API key XORed with the webhook secret, so whoever learns either reads the other wherever the two overlap.
- [ ] The pad's first 25 bytes
  > The pad cancels out of the XOR; pad bytes show only to someone who also knows a secret.

## Six hex digits make 16.7 million codes. After about how many keys do two probably share one?

- [ ] About 16.7 million, one key for every code
  > That is about how many tries it takes to match one given code; any two match far sooner.
- [ ] About 8.4 million, half of the codes
  > Long before half the codes are used, some pair of keys has almost surely matched.
- [x] About 4,800, near the square root of 16.7 million
  > Every new key can match any earlier one, so by about 4,800 keys the odds of a shared code pass even: the birthday bound.

## What does a pad need to hide a message completely?

- [x] Random bytes as many as the message's, used for that message alone
  > Then every message of that length is equally likely, whatever the sealed bytes are, and the pad gives away nothing but the length.
- [ ] Bytes kept safe at both ends, reused only for Apogee's own messages
  > Reuse is the flaw whoever sends: any two messages under the same pad bytes leak their XOR.
- [ ] A pad of 32 bytes, as long as a key
  > A 32-byte pad can't cover a 43-byte API key without using some bytes twice: it must be as long as the message.
