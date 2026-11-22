Give the channel room for one value, so the goroutine can always send its result and finish, even after `run` has returned.

A `select` with two cases waits for whichever is ready first: a value on the channel means `Serve` stopped by itself, and `<-ctx.Done()` means it is time to stop.

In the `ctx.Done()` case, return what `stop(srv, t.Shutdown)` returns; `Shutdown` makes `Serve` return too, and the buffered channel takes its result. If the test says `run` returned while `ctx` still lived, nothing is waiting in a `select`.
