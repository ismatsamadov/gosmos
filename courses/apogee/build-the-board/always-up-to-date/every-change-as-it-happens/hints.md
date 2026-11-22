Replace the receive from `r.Context().Done()` at the end of `serveEvents` with a `for` loop around a `select` of two cases.

One case receives a view from `b.updates` and passes it to `writeEvent`; when that returns an error, the screen is gone, so return. The other case is `<-r.Context().Done()`, and returns.

The view on `b.updates` already carries its version, so `writeEvent` gives it the right `id`. Without the context case the loop would wait on `b.updates` after the screen left, and the stream would never end.
