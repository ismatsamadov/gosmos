## In `GET /recordings HTTP/1.1`, what is `/recordings`?

- [x] The target: the path, exactly as the client wrote it
  > The station server takes it as sent and decides what it names; nothing is cleaned up or guessed.
- [ ] The name of the host the request is for
  > That is the Host field, a header line of its own; this target holds only the path.
- [ ] A file the server reads from its disk
  > The target names something the server knows, here the list of recordings, and the server decides how to answer it.

## Why is `GET /recordings HTTP/1.0` answered 505 and not 400?

- [ ] Because every refused request line is a 505
  > A line of the wrong shape gets 400; 505 is kept for one case only.
- [x] Its shape is right, and only its version is one the server doesn't speak
  > A 400 says the line itself is crooked; the station server's 505 says the line is good but its version is not HTTP/1.1.
- [ ] Because an HTTP/1.0 request has no request line
  > It has one, of exactly this shape; only the version differs.

## What becomes of `get /recordings HTTP/1.1`?

- [ ] The server reads it as GET, since methods ignore letter case
  > Method names are case-sensitive, so `get` is a different method from GET.
- [ ] parseRequestLine refuses it with 400 bad request line
  > Its shape is right, three parts split by single spaces, so it passes; the method is judged later.
- [x] It passes as `get`, and a later check refuses it, since `get` is not GET
  > Methods are case-sensitive and the shape is right, so parseRequestLine returns `get` as sent; on Friday, readRequest answers any method but GET, HEAD and PUT with 501.
