`credentials(r, "Bearer")` gives the token or a 401 error; only when it worked, hand the token to `verifyToken` with `key` and `now()`. The customer ID it returns must be a key of `customers`, or the error is `unauthorized("invalid token")`.

Carry one error through the checks and test it once: if it isn't nil, set the header on `w.Header()`, call `writeError`, then `return` before anything else runs.

If the tests say `next` ran after a 401, you are missing that `return`. If `customerFrom` finds nothing, you called `next` with `r` itself: pass it `r.WithContext(withCustomer(r.Context(), c))`.
