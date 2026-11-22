## What lets a client send another request on the same connection?

- [ ] Having read the answer's head
  > The body still follows the head; until the client has read it to its end, its bytes would be taken for the next answer.
- [ ] A `Connection: keep-alive` field in the answer
  > HTTP/1.1 keeps connections open without it, and an answer with no length or chunks couldn't share the line anyway.
- [x] Knowing where each answer ends, by its length or its last chunk
  > Then no byte of one answer can be taken for the next, so both ends can use the line again.

## A PUT gets 413 before its body is read. Why does the server hang up then?

- [x] The unread body is still on the line, and would be read as the next request
  > Its bytes would be taken for the next request, and reading them all just to skip them is what the 413 refused to do.
- [ ] A 413 says the client is broken, so it gets no second chance
  > A 413 says only that this body is too big; the client may send a smaller one on a new connection.
- [ ] Go's client sends `Connection: close` with every PUT
  > It doesn't; the server decides to close, and says so in its answer.

## The feed stops after two lines, with no last chunk. What does that tell the collector?

- [ ] The server had nothing more to send
  > A feed that ends on purpose always sends its last chunk, as `0\r\n\r\n`.
- [x] The stream was cut, so it reconnects and reads `/recordings` for what it missed
  > Only the zero chunk tells a whole stream from a cut one; without it, the collector can't trust that nothing came after.
- [ ] Two recordings were stored, and each feed carries two
  > A feed carries every recording stored while it is open, however many that is.
