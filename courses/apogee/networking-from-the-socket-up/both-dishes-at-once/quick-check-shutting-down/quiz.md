## What does the accept loop's blocked `Accept` do when `s.ln.Close()` runs?

- [x] It returns at once, with an error matching `net.ErrClosed`
  > Closing the listener wakes the `Accept` waiting on it, and `serve` takes that error as its signal to return.
- [ ] It waits for one more watcher, then fails
  > A closed listener accepts nobody: the waiting `Accept` returns right away.
- [ ] It returns a connection that is already closed
  > After its listener closes, `Accept` returns no connection at all, only the error.

## A watcher's goroutine is blocked in `Read`. What makes that `Read` return?

- [ ] Closing the listener the watcher connected through
  > The listener only takes new connections; the ones it accepted live on without it.
- [ ] Calling `s.wg.Wait()`
  > `Wait` only waits for goroutines to end; it touches no connection.
- [x] Closing that watcher's connection
  > A `Read` blocked on a connection returns with an error as soon as the connection closes, so its goroutine can end.

## Why does this draft never return while three watchers are connected?

- [ ] `s.ln.Close()` blocks until every watcher hangs up
  > Closing a listener returns at once and leaves the accepted connections alone.
- [x] `Wait` waits for goroutines that end only when their connections close, after it
  > Every close must come before `Wait`: the goroutines it waits for are blocked on connections nobody has closed yet.
- [ ] `CloseAll` waits for each watcher to read its last line
  > `CloseAll` closes the connections and returns; it waits for nobody.
