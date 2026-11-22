## Someone on the island's line changes a POINT's azimuth and computes a new CRC for it. What does this check do?

- [ ] Refuses it, because nobody can compute a CRC without a key
  > A CRC-32 has no key: anyone who has the bytes can compute it.
- [x] Accepts it: a CRC catches accidents, and only a signature stops a deliberate change
  > Anyone can recompute a CRC. A signature needs a key, which is why Pipit-2's commands carry one.
- [ ] Refuses it, because the new CRC no longer matches the frame's length
  > The CRC never covers the length, and a recomputed CRC matches the changed bytes.

## Why does the check come before any field is read?

- [x] A damaged frame's fields can't be trusted, so nothing should act on them
  > Checked first, a damaged STATUS never puts a wrong azimuth on anyone's screen.
- [ ] Reading the fields first would change the bytes the CRC covers
  > Reading changes nothing; the point is that no damaged field is ever used.
- [ ] The CRC comes first in the payload
  > It comes last: the payload's final four bytes.

## Is `crc32.ChecksumIEEE` the CRC-32 that downlinked files carry?

- [ ] No: frames carry a shorter CRC-16
  > A frame's CRC is four bytes, a CRC-32, like a file's.
- [ ] No: the frames' CRC is written little-endian
  > Byte order is how the number is written down, not which CRC it is; frames write theirs big-endian.
- [x] Yes: the same IEEE CRC-32
  > A downlinked file's CRC, sent after its chunks, is the CRC-32 (IEEE) that `crc32.ChecksumIEEE` computes.
