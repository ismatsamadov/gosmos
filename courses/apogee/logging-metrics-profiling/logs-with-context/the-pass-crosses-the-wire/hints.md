In `RoundTrip`, ask `passFrom` about the request's context. Only when it holds a pass, make a clone of `r` with `Clone`, set `Apogee-Pass` on the clone's header, and send the clone through `t.base` instead of `r`.

`acceptPass` returns an `http.HandlerFunc` that reads the request's `Apogee-Pass` header. When it isn't empty, `next` gets the request made by `WithContext` from a context `withPass` made; otherwise `next` gets the request as it came.

If the test says the caller's request now has the header, the header was set on `r` itself rather than on its clone; if the station got no header, check that the clone, not `r`, is what goes to `t.base`.
