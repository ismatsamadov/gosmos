Build it as `noStore` is built: return an `http.HandlerFunc` that sets the four headers on `w.Header()` and then calls `next.ServeHTTP(w, r)`.

Copy each value exactly, the single quotes inside the policy included: `'self'` and `'none'` are written with them.

If the redirect or the 403 lacks the headers, check the order: a header set after `next` returns comes too late, because the answer was written inside that call.
