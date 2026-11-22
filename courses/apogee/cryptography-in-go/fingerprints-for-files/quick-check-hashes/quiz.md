## Anyone can compute a SHA-256 digest. When does one prove Cirrine's file unchanged?

- [ ] Always, since nobody can make two files with one digest
  > Nobody can, but whoever changes a file can send the changed file's own digest with it.
- [x] Only when the digest itself reaches Cirrine safely
  > A changed file with its own fresh digest checks out perfectly, so the digest has to arrive by a path nobody can change.
- [ ] Only for files small enough to hash at once
  > Hashing a stream gives the same digest as hashing it whole, for a file of any size.

## How did Mirela's forgery keep the report's CRC-32?

- [ ] She tried bytes until one set happened to match
  > No trying was needed: a CRC's arithmetic runs backward, so the four bytes can be worked out directly.
- [ ] A CRC-32 ignores the order of the digits
  > It doesn't: swapping 7 and 1 changed the CRC-32, and the four added bytes changed it back.
- [x] Four bytes she worked out bring any file back to a chosen CRC-32
  > A CRC catches noise, not intent: anyone can compute the four bytes that restore it.

## How long is the SHA-256 digest of an empty file, and of a 10.9 GB image?

- [x] 32 bytes each
  > SHA-256 always gives 32 bytes, from no bytes at all or from 10.9 GB.
- [ ] 0 bytes, and 32 bytes
  > An empty file still has a digest, `e3b0c442…b855`, 32 bytes like any other.
- [ ] 32 bytes, and more for a file that large
  > The digest never grows with the file; that's what lets it stand for any file.
