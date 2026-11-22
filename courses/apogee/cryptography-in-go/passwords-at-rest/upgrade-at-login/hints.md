At the top of `checkStored`, `strings.CutPrefix` with `"sha256$"` gives you the old record's 64 digits and whether it was an old record at all. Only a record that isn't old goes on to `parseRecord`.

For an old record, hash the password with `sha256.Sum256`, write its digest as hex, and pass both texts to `subtle.ConstantTimeCompare` as byte slices. An old record is always weak; a parsed one is weak when `r.Iterations` is below the constant `iterations`.

In `login`, keep the `weak` result instead of dropping it. Return early, with no new record, on an error, a wrong password or a strong record; otherwise call `newRecord` and return what it gives you.
