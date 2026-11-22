In `RoundTrip`, `spanFrom(r.Context())` tells you whether there is a span to send. Only then, replace `r` with `r.Clone(r.Context())` and set the header on the clone, as the trace ID, a `-`, and the span ID.

`continueTrace` returns an `http.HandlerFunc`. Inside it, start from `r.Context()`, put the parent in with `withSpan` only when `remoteParent` returns true, call `t.Start` with the result, defer `t.End`, and serve `next` with `r.WithContext` and the context `Start` returned.

If the station's span has no parent although a header came, the span was started from `r.Context()` itself rather than from the context holding the remote parent. If headquarters' own request carries the header afterwards, it was set before cloning.
