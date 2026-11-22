In `WriteHeader`, store `code` in `s.Status` only while `s.Status` is still 0, then call `s.ResponseWriter.WriteHeader` so the real writer gets it too.

In `Write`, a `Status` of 0 means no status came first, and an answer written that way is a 200. Pass `p` to `s.ResponseWriter.Write`, then add the count it returns to `s.Bytes` and return that count and its error.

`Unwrap` takes no arguments and returns the embedded field, which you reach as `s.ResponseWriter`. If `Flush` still fails, compare your method with the exact signature `Unwrap() http.ResponseWriter`.
