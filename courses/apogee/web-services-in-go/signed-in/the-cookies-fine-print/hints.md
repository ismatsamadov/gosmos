Each attribute is a field of `http.Cookie`: add `Path`, `MaxAge`, `HttpOnly`, `Secure` and `SameSite` to the literal that `sessionCookie` returns, and keep `Name` and `Value` as they are.

`sessionLife` is a `time.Duration`. Its `Seconds` method gives the seconds as a `float64`, and `MaxAge` wants an `int`, so convert it.

If the header line still differs, read it against the expected one attribute by attribute: Go always writes them in the same order, so the missing or wrong one stands out.
