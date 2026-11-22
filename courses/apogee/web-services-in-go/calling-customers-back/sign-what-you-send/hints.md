`hmac.New(sha256.New, secret)` returns a hash you `Write` to: first `ts` in decimal followed by a dot, made with `strconv.FormatInt`, then `body`. `Sum(nil)` returns 32 bytes, and `hex.EncodeToString` turns them into 64 digits.

In `post`, `body` already holds the exact bytes the request sends. Set both headers on `req` before `c.Do`, the timestamp as `strconv.FormatInt(ts, 10)`.

If the receiver's computation differs from yours, check the order: the timestamp, a dot, then the body, with nothing between them. If only the case a second later fails, the timestamp isn't in what you sign.
