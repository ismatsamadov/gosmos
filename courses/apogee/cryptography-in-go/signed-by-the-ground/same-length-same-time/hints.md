`sha256.Sum256` takes a `[]byte`, so convert each phrase, and it returns an array: slice it with `[:]` to hand it to `subtle.ConstantTimeCompare`.

`subtle.ConstantTimeCompare` returns an `int`, 1 when the slices are equal and 0 otherwise, so compare its result with 1 to get the `bool` that `samePhrase` returns.

If `TestNoEquals` fails, look for any `==` or `!=` between the phrases, their lengths or their digests: the only comparison left may be the one between `ConstantTimeCompare`'s result and 1.
