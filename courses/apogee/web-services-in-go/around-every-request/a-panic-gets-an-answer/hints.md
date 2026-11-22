Put the `defer` first in the handler, before `next.ServeHTTP`, and call `recover()` inside the deferred function. A nil result means nothing panicked, and there is nothing to do.

Compare the recovered value with `http.ErrAbortHandler` using `==`, and when they match, call `panic` with the value again before you log or write anything.

For any other value, log it with `l.Printf` after the ID from `requestIDFrom(r.Context())`, then give `writeError` an error that isn't an `*APIError`, such as one made with `fmt.Errorf`: it answers 500 `internal` and keeps the panic's text out of the answer.
