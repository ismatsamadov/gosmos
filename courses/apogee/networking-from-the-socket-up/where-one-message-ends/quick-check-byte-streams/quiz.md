## dishd wrote the PING and the QUERY with two Writes. How can the controller receive them?

- [ ] Always as two Reads of 11 bytes, one per Write
  > A connection keeps no edges between Writes: the 22 bytes can come in any number of pieces.
- [x] In any pieces: all 22 bytes in one Read, or one frame split over several
  > A connection carries a stream of bytes, and a reader sees only what has arrived, cut wherever the line cut it.
- [ ] As one Read of 22 bytes, because they were sent one after the other
  > Bytes written together can still arrive in pieces, and bytes written apart can arrive together.

## `buf` is 11 bytes long and the first Read inside `io.ReadFull(conn, buf)` brings 3. What happens next?

- [ ] It returns 3 and nil, as a single `Read` would
  > That is a single `Read`'s answer; `ReadFull` keeps going until `buf` is full.
- [ ] It returns `io.ErrUnexpectedEOF` at once
  > That error is for a stream that ends before `buf` is full; 3 bytes is only a slow line.
- [x] It reads again until `buf` holds 11 bytes, or the stream fails
  > `ReadFull` calls `Read` as often as it takes, and returns an error only when the stream ends or fails first.

## Why does every frame start with its length, when the banner ends with a newline?

- [x] A frame can hold any byte, `0x0a` included, so no byte can mark where it ends
  > Binary data has no byte to spare for a delimiter: a WELCOME after command 10 has a newline in its middle.
- [ ] Two bytes of length take less room than a newline
  > Two bytes are more than one; the length is there because a delimiter can't work for frames.
- [ ] The protocol keeps the newline for the banner alone
  > Nothing keeps it: a frame may contain `0x0a` anywhere, which is why a delimiter can't find its end.
