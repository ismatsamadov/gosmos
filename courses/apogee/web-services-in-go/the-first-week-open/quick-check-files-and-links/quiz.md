## A resume asks for `Range: bytes=40000-` and gets `200 OK`. What must the client do?

- [ ] Add the body to the 40,000 bytes it already has
  > A 200 carries the whole file from its first byte, so adding it would repeat the first 40,000 bytes.
- [x] Throw away its 40,000 bytes and keep this body, the whole file
  > Only a 206 continues a file. A 200 means the server sent everything, perhaps because the file changed since.
- [ ] Ask again with a smaller range until a 206 comes
  > A server that answers a range with 200 isn't serving ranges of this file; asking again changes nothing.

## Why does `readCommands` wrap `r.Body` in `http.MaxBytesReader` before it reads the form?

- [x] Reading stops at 64 KiB, so a 1.9 GB upload costs the service 64 KiB
  > The read that crosses the limit fails with a `*http.MaxBytesError`, which the service answers with 413.
- [ ] It makes `NextPart` skip any part larger than the limit
  > It skips nothing: the read that crosses the limit fails, and the whole form with it.
- [ ] It tells the client the limit before the upload starts
  > The client learns of the limit only from the 413 that ends its upload.

## Why does a link's signature cover its expiry as well as its path?

- [ ] So that two links to the same file never look alike
  > Two links to one file made in the same second are alike, and that harms nobody.
- [ ] Because `hmac.Equal` compares only values of the same length
  > `hmac.Equal` compares any two byte slices; signing the expiry is about what a link may claim.
- [x] Otherwise anyone holding a link could move its expiry on and keep it working
  > With the expiry signed, a changed `expires` breaks the signature, and the link is refused as invalid.
