## A's call runs the load. What does B's call of `get` return?

- [ ] `nil`, since the load hasn't finished when B asks
  > B's `Do` doesn't return early: it waits until A's function has returned.
- [x] The network A's load made: B's `Do` waits for it and sees what it wrote
  > When `once.Do` returns, every write its function made is visible to the caller, even one whose call didn't run it.
- [ ] A network of its own, from a second load
  > `Do` runs its function only once, however many goroutines call it.

## A function from `sync.OnceValues` got an error from its first load. What do later calls return?

- [x] The same error, every time, without loading again
  > `OnceValues` keeps its first results, error included: a load that may fail and later succeed needs something else.
- [ ] A new try at the load, until one succeeds
  > Nothing calls the load a second time, whatever it returned.
- [ ] The network, once the file server is back
  > The first results are kept for good, so a later recovery never reaches the callers.

## What happens when the function given to `once.Do` calls the same `once.Do`?

- [ ] It runs a second time, inside the first
  > `Do` never runs its function twice.
- [ ] The inner call returns at once and does nothing
  > The inner call waits for the first call to finish, as every other caller does.
- [x] The inner call waits forever for the outer one, which waits for it
  > The outer `Do` can't return until its function does, and that function is waiting in the inner `Do`: a deadlock.
