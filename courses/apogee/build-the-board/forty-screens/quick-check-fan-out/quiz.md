## Two streams receive from the same channel. Who gets a view sent on it?

- [ ] Both streams, one copy each
  > A channel delivers each value sent on it to exactly one receiver, never to both.
- [x] Exactly one of them, whichever receives first
  > A sent value reaches one receiver, so the other screen misses that change; that is why each stream needs its own channel.
- [ ] Neither, until both are ready to receive
  > A buffered send needs no receiver at all, and a waiting receiver takes the value at once.

## The channel's 4 slots are full. What does this `select` do?

- [ ] It waits until one of the streams receives a view
  > A plain send would wait; with `default` the `select` never does.
- [ ] It panics, because the channel is full
  > A full channel never panics: a send on it waits, or `select` takes another case.
- [x] It runs `default` at once, and the view is not sent
  > When no case is ready, `default` runs, which is how the hub notices a stream that has fallen behind.

## Why does only the hub close a stream's channel?

- [x] It knows whether the channel is already closed, and closing one twice panics
  > The hub closes a channel only while it still holds it, under its mutex, so each channel is closed exactly once.
- [ ] Only the goroutine that made a channel may close it
  > Any goroutine that holds the channel's sending end may close it; the danger is closing it twice.
- [ ] A closed channel stops the server from shutting down
  > Closed channels are what lets streams end, which helps the server shut down.
