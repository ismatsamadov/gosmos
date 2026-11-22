`strconv.FormatInt(expires.Unix(), 10)` writes the expiry in Unix seconds. The customer, a dot and that number are the claims: the token starts with them, and they are what you sign.

Make the HMAC with `hmac.New(sha256.New, key)`, write the claims to it as bytes, and hex-encode what its `Sum(nil)` returns.

If the signature differs from the one the test wants, check exactly what you signed: the customer, a dot and the expiry, the very text the token starts with, and not the customer alone. The token is that text, a dot, then the 64 hex digits.
