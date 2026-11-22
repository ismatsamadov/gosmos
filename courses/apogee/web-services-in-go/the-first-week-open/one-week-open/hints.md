`s.private` takes a handler function such as `f.upload` and returns the handler to register with `mux.Handle`: the token check, then the customer's limit.

A `metered` keeps the store it wraps in its `Store` field and the `*Quota` it counts against in `quota`. The booking handler reads `s.Store` on every request, so replacing that field is all it takes.

The notifier delivers through `app.Notifier.Sender.Client`. If the weather service still gets Brinewatch's event, that field still holds the client `newApp` made.
