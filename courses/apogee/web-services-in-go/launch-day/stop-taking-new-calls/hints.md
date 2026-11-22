`context.WithTimeout` makes a context that ends after `within`, and a `cancel` function to call, with `defer`, once you are done with it.

`Shutdown` returns nil once every request in flight has finished, or the context's error if the context ended first: return what it returns.

Make the timeout only when `within` is more than 0, and otherwise pass `context.Background()`, which never ends. If the case with no limit fails with `context deadline exceeded`, you gave `WithTimeout` a timeout of 0.
