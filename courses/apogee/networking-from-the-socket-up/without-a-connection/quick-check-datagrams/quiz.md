## What does the second `ReadFrom` put in `buf`?

- [ ] The last 8 bytes of `a`, which didn't fit the first time
  > That is how a TCP stream behaves. A datagram is read whole or cut, and the part of `a` that didn't fit is gone.
- [ ] Nothing: it fails, because `a` was cut
  > A cut datagram brings no error at all, to the first read or the next one.
- [x] The first 64 bytes of `b`
  > Each `ReadFrom` takes exactly one datagram: the second takes `b` and cuts it to the 64 bytes of `buf`.

## The first `WriteTo` returned 72 and nil. What does the recorder know?

- [x] Only that the datagram left its socket
  > Nothing comes back to say a datagram arrived; one sent where nobody listens vanishes without an error.
- [ ] That the live view has received it
  > No answer or acknowledgment ever comes back on UDP unless the program at the other end sends one.
- [ ] That a socket is listening at that address
  > With no connection there is no handshake, so nothing checks that anyone is there.

## Why does `pack` keep every datagram to whole lines?

- [ ] Because `ReadFrom` stops at the end of a line
  > `ReadFrom` knows nothing of lines: it returns a datagram, whatever bytes it holds.
- [x] So a lost datagram costs only its own frames
  > A line split across two datagrams would be lost if either one were, and half a line is no frame at all.
- [ ] So the live view can find each datagram's end by its newline
  > A datagram's end comes with it, and the last line in one has no newline after it.
