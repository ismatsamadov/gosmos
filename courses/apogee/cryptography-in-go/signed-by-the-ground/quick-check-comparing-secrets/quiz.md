## Why does the uplink compare a command's tag with `hmac.Equal`?

- [ ] `bytes.Equal` sometimes gets the answer wrong
  > It gets every answer right; only the time it takes gives too much away.
- [x] A quick refusal would show how much of a forged tag was right
  > Stopping at the first difference lets an attacker find a tag piece by piece; `hmac.Equal` takes the same time wherever the tags differ.
- [ ] `hmac.Equal` also checks which key made the tag
  > It compares bytes only; the key went into making the tag it compares against.

## Cirrine's archive compares a file's digest with the manifest's. Must that take the same time whatever the bytes?

- [x] No: both digests are public, so its timing gives nothing away
  > Anyone can compute a file's digest, so there is nothing secret for the timing to reveal.
- [ ] Yes: every comparison of digests must take the same time
  > Equal timing protects secrets, and a public digest isn't one.
- [ ] Yes, because a 10.9 GB file takes longer to hash
  > Hashing time depends on the file's size, which the manifest already shows.

## Why does `samePhrase` hash both phrases before `subtle.ConstantTimeCompare`?

- [ ] Hashing makes the comparison faster
  > Hashing adds work; it's done for the lengths, not for speed.
- [ ] `subtle.ConstantTimeCompare` accepts only 32-byte slices
  > It takes slices of any length, which is the trouble: it answers at once when their lengths differ.
- [x] It answers at once for different lengths, and digests are always 32 bytes
  > Every digest is 32 bytes, so the comparison's timing never gives away how long the stored phrase is.
