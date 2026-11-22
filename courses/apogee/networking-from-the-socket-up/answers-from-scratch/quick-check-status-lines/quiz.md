## Which part of `HTTP/1.1 413 Request Entity Too Large` does a client act on?

- [x] The code, 413
  > A client decides by the number; the reason phrase is for people, and servers word it differently.
- [ ] The reason phrase, read word for word
  > Newer servers say `Content Too Large` for the same 413, so no client depends on the words.
- [ ] The version, `HTTP/1.1`
  > The version says how to read the rest of the answer; the code says what happened.

## Where does the head of an answer end?

- [ ] After its `Content-Length` field
  > Fields come in whatever order the server writes them; the station server puts `Location` after it.
- [ ] At the CR LF that ends the status line
  > That ends only the first line; the fields follow, one per line.
- [x] At the blank line, a CR LF on its own
  > Everything before it is the head and everything after it the body, counted by `Content-Length` or by chunks.

## An answer gives neither `Content-Length` nor chunks. Where does its body end?

- [ ] At the next blank line
  > A body can hold blank lines: a recording's bytes can be anything, CR and LF included.
- [x] Where the connection closes, so a cut answer looks whole
  > With no length and no last chunk, closing is the only end there is, and a line that drops halfway ends the body just the same.
- [ ] Where no byte arrives for a second
  > A slow line pauses too, and a pause can't tell the end of a body from a delay.
