Derive exactly as the record was made: the same hash, `pw`, the record's salt and its iteration count, and a key as long as `r.Key`, with `"crypto/pbkdf2"` added to the imports. If `pbkdf2.Key` returns an error, return `false` and that error.

`subtle.ConstantTimeCompare` returns an `int`, not a `bool`: 1 when the two slices are equal and 0 otherwise. Compare its result with 1.

If the bench's case fails while Mirela's pass, the count is fixed at 600,000 instead of read from `r.Iterations`. If `TestConstantTime` fails, its message names the line where two keys are compared some other way.
