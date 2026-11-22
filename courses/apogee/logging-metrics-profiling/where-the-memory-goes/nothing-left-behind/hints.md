Start `runPass` with the `select`, which takes a worker or gives up: one case sends `struct{}{}` on `u.workers` and does nothing more, the other receives from `ctx.Done()` and returns the wrapped `ErrNoWorker`. After it, the worker is yours until `runPass` returns, which is what the `defer` is for.

Then call `u.fly` right in `runPass`, with `ctx` where the starter passed `context.Background()`. The goroutine, the `done` channel and the last `select` all go: the radio now stops by itself when the pass's context ends.

If the test reads `context deadline exceeded`, the old `select` on `ctx.Done()` is still answering for the pass; if it counts goroutines, a `go` statement is still there. Wrap both errors with `%w`, as `pass <ID>: ...`, so `errors.Is` still finds them.
