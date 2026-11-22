After the lock is released, call `s.closeConns(false)`, which takes the lock itself and closes only the connections waiting between requests. The busy ones finish their answer and hang up, because `send` sees the server stopping. `s.Store.Close()` then ends every feed with its last chunk.

`s.wg.Wait()` takes no context, so run it in a goroutine that closes a `done` channel when it returns, and `select` on `done` and `ctx.Done()`.

When `ctx` ends first, `s.closeConns(true)` cuts what is left; wait on `done` again, so nothing outlives `Shutdown`, then return `ctx.Err()`. If the upload in flight is cut, every connection was closed at once; if the waiting one stays open, `closeConns(false)` never ran.
