The `switch` already has a `case` for `sha256`. Add one for the exact name `sha3-256` that returns a new SHA3-256 hash and no error.

`sha3.New256()` returns a value that satisfies `hash.Hash`, so it can be returned as it is. Remove the `_ = sha3.New256` line once you use it.

Leave the last `return` as it is: every other name, `SHA3-256` in capitals included, must reach it, so that `fingerprintAs` stops before it reads anything.
