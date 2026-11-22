## When `readBanner` returns, where can the first frame's bytes be?

- [x] Already in `br`'s buffer, read from the connection along with the banner
  > `br` reads in big gulps, so bytes that arrived with the banner wait in its buffer, gone from the connection.
- [ ] Only in the connection, because `ReadString` never reads past the newline
  > `ReadString` returns at the newline, but `br` may have read further, so its buffer can already hold what came after it.
- [ ] Nowhere: `ReadString` throws away the rest of what it read
  > Nothing is thrown away: what `br` read past the line stays in its buffer for the next read.

## Which reader must read every frame after the banner?

- [ ] A new `bufio.Reader` per frame, so that frames never mix
  > A new reader can't see the bytes the first one buffered, and frames go missing that way.
- [ ] `conn` itself, since frames are binary
  > Reading `conn` directly skips whatever `br` has already buffered.
- [x] `br`, for the rest of the connection's life
  > Once a connection is wrapped, every read goes through the same reader, which holds whatever it read ahead.

## A frame's length reads `00 05`. What does `readPayload` do?

- [ ] Reads the 5 bytes, then finds them too few for a frame
  > Checking after reading trusts a length that may be damaged: `ff ff` would wait for 65,535 bytes.
- [x] Returns `ErrBadLength` before reading anything more
  > A length under 9 can't be a frame, so it is refused at once, before any byte is read or any slice made.
- [ ] Reads on until the CRC shows where the frame ends
  > A CRC can't mark an end: only the length says how many bytes to read.
