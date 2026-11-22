`credentials(r, "Key")` gives the key, or the 401 to answer with. Hand the key to `checkKey` with `s.Customers`; when it says no, the error to answer with is `unauthorized("invalid api key")`.

Set `WWW-Authenticate` on `w.Header()` before calling `writeError`: a header set after the answer is written never reaches the caller.

Add an hour to `s.Now()` and truncate that to the second once, then use the same time for `signToken` and for `ExpiresAt`, and answer with `writeJSON` and 200. If `expires_at` ends in `.6`, you truncated only one of them.
