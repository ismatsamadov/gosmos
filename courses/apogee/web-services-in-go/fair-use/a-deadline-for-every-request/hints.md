Return an `http.HandlerFunc` as other middleware does. Inside it, when `d` is zero or less, call `next` with `r` as it is and return.

`context.WithTimeout` takes the request's own context and `d`, and returns a new context and its `cancel` function: defer the cancel, then pass `next` the request that `r.WithContext` makes from the new context. In `writeError`, two more cases with `errors.Is` find the context errors, however the store wrapped them.

If the stalled store answers 500 after half a second, `next` got the old request, whose context has no deadline. If the test says the context is still live after the handler returned, the `cancel` never ran.
