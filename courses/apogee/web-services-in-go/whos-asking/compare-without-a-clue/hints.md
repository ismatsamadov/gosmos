`strings.Cut(key, ".")` returns what comes before the first dot: the customer ID. Look it up in `keys`; when the table has no such customer, the key can't be right.

`subtle.ConstantTimeCompare` takes two `[]byte`, so convert the key you were given and the stored one with `[]byte(...)`. It returns an `int`, 1 for equal, not a `bool`.

If the tests say `checkKey` never calls `subtle.ConstantTimeCompare`, make the call inside `checkKey` itself, in place of `==`. If a key one digit short gets through, you are comparing only part of it: compare the whole key with the whole stored one.
