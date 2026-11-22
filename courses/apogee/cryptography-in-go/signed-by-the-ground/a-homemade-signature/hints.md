`hmac.New(sha256.New, key)` gives a hash with the key built in. Write the command into it, then `Sum(nil)` returns the 32-byte tag.

Drop the `sha256.Sum256` line and the `append` calls that put the key in front of the command: with HMAC, the key never goes into the message at all. Remove the `_ = hmac.New` line too.

If `TestTag` fails while `TestLengthExtension` passes, check that the key goes to `hmac.New` and only the command is written; the lab bench's signature since 2027, `512545bc…`, is the HMAC-SHA256 of the command alone.
