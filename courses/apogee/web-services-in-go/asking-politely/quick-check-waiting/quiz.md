## How long does `http.Get` wait for a server that has stalled?

- [ ] Two minutes, the default client's limit
  > `http.DefaultClient` has no `Timeout` at all, so it sets no limit.
- [x] Forever: it uses `http.DefaultClient`, which has no `Timeout`
  > With no limit set, only the server, or the end of the program, ends the wait.
- [ ] Thirty seconds, after which the server gives up
  > Nothing makes a stalled server give up; that is why the client needs a limit.

## Does a client's `Timeout` still count while the body is being read?

- [x] Yes: it covers the whole call, reading the body included
  > The timer keeps running after `Do` returns, and it stops a slow body too.
- [ ] No: it stops counting once the headers have arrived
  > The headers are only part of the call; the limit covers the body as well.
- [ ] Only when the body is read with `io.ReadAll`
  > How the body is read makes no difference: the same timer covers it.

## In the code above, when does `c.Do` return?

- [ ] After 2 seconds, the client's `Timeout`
  > The context ends first, at half a second, and that ends the call.
- [ ] After 2.5 seconds, the two limits added together
  > The limits don't add up: the first one to end stops the call.
- [x] After half a second, when the context ends
  > Whichever ends first, the context or the client's `Timeout`, ends the call.
