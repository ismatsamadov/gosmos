Return an `http.HandlerFunc` from the inner function, and add `Vary: Origin` before anything else, so that answers to other origins and to programs carry it too.

Compare the request's `Origin` header with `origin` using `==`. A prefix or a part of it would also trust `https://plan.furrowlight.example.elsewhere.example`.

For the trusted origin's `OPTIONS`, set the headers, write 204 and return without calling `next`: the API behind it refuses a preflight, which never carries a token.
