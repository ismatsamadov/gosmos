## A request says `Content-Length: 4096`. Where does its body end?

- [x] After exactly 4,096 bytes, whatever they hold
  > The length counts bytes, and zeros, CRs and LFs are bytes like any other, so nothing inside the body can end it sooner.
- [ ] At the first CR LF CR LF after the head
  > A blank line ends the head; inside a body those four bytes are data, and a recording is full of them.
- [ ] When the client closes its side of the connection
  > The length says where the body ends, so the connection stays open for the answer; a stream that ends sooner is a cut body.

## Where does a chunked body end?

- [ ] At the first chunk smaller than the one before it
  > Chunks come in any sizes, in any order; each size line says only how long its own chunk is.
- [ ] When the client stops sending for a while
  > A pause is a slow line, not an end: the server's deadline turns a long one into 408, never into a whole body.
- [x] At the chunk of size 0 and the blank line after its trailer
  > A size of 0 marks the last chunk, any trailer fields follow it, and the blank line after them ends the body.

## Why does the recorder send each pass chunked?

- [ ] Chunked bodies are smaller on the wire
  > Each chunk adds its size line and a CR LF, so a chunked body is a few bytes longer than the same bytes with a length.
- [x] It can start before it knows how long the pass will be
  > A `Content-Length` must come before the body, and the recorder knows a pass's length only at LOS, when the satellite sets.
- [ ] The station server takes recordings only in chunks
  > It takes either; a test recording of known size comes with a `Content-Length`.
