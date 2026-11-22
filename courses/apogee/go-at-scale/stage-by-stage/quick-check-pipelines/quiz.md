## The Board stops reading. What happens to a stage whose send has no `ctx.Done()` case?

- [x] Its goroutine waits on that send forever, even after the run is canceled
  > Nothing interrupts a plain send: only a `select` with `<-ctx.Done()` lets a cancel reach the goroutine.
- [ ] Go ends the goroutine once the Board's goroutine has returned
  > Go never ends a blocked goroutine for you; it waits there until the program exits.
- [ ] The send gives up with `context.Canceled` once the run is canceled
  > A channel send knows nothing about contexts: it waits until a receiver takes the value.

## `readFeed` cancels the run with `cancel(err)`. What do `ctx.Err()` and `context.Cause(ctx)` return?

- [ ] Both return `err`
  > `ctx.Err()` only ever says why a context ended in general, `context.Canceled` here; the error given to cancel comes from `Cause`.
- [x] `context.Canceled`, then `err`
  > `Err` says that the context ended; `Cause` says why, with the error given to the cancel function.
- [ ] `err`, then `nil`
  > The roles are the other way round, and the cause of a canceled context is never `nil`.

## When may a fan-in close its output?

- [ ] As soon as its first input closes
  > The other forwarders may still send, and a send on a closed channel panics.
- [ ] The moment the context is canceled
  > A forwarder may be about to send; the output closes only after every forwarder has returned.
- [x] Once every forwarding goroutine has returned
  > Then nothing can send on it any more, so closing it is safe: `wg.Wait()`, then `close`.
